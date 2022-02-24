package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"

	//"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"
	"v/domain"
	"v/ntlmssp"
)

type Token struct {
	tokenString string
	mx          sync.RWMutex
}

type Client struct {
	http.Client
	token             Token
	username          string
	password          string
	tokenURL          string
	getProcessesURL   string
	getProcessURL     string
	getIncidentsURL   string
	retryTaskURL      string
	retryJobURL       string
	activitySearchURL string
	modificationURL   string
	updateVarsURL     string
	managerRoleURL    string
}

// NewClient returns new Client
func NewClient() *Client {
	c := &Client{
		username:          os.Getenv("USERNAME"),
		password:          os.Getenv("PASSWORD"),
		tokenURL:          os.Getenv("TOKEN_URL"),
		getProcessesURL:   os.Getenv("PROCESSES_URL"),
		getProcessURL:     os.Getenv("PROCESS_URL"),
		getIncidentsURL:   os.Getenv("GET_INCIDENT_URL"),
		retryTaskURL:      os.Getenv("RETRY_TASK_URL"),
		retryJobURL:       os.Getenv("RESTRY_JOB_URL"),
		activitySearchURL: os.Getenv("ACTIVITY_SEARCH_URL"),
		modificationURL:   os.Getenv("MODIFICATION_URL"),
		updateVarsURL:     os.Getenv("UPDATE_VARS_URL"),
	}
	if err := c.setToken(); err != nil {
		c.setToken()
	}
	c.Transport = ntlmssp.Negotiator{
		RoundTripper: &http.Transport{
			MaxIdleConnsPerHost: 25,
			DialContext: (&net.Dialer{
				Timeout: time.Second * 5,
			}).DialContext,
		},
	}
	token := c.getToken()
	if token == "" {
		c.setToken()
	}
	return c

}

// SetToken allows user set their own token
func (c *Client) SetToken(token string) error {
	if token == "" {
		return domain.ErrTokenNotFound
	}
	c.token.tokenString = token // TODO set their token if any!
	return nil
}

// setToken gets a token for given username and password and sets it as token for Client.token
func (c *Client) setToken() error {
	if c.tokenURL == "" {
		return domain.ErrTokenURLNotFound
	}

	req, err := http.NewRequest(http.MethodGet, c.tokenURL, nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(c.username, c.password)
	res, err := c.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	var dst bytes.Buffer
	if _, err := io.Copy(&dst, res.Body); err != nil {
		return err
	}
	c.token.mx.Lock()
	defer c.token.mx.Unlock()
	var token string
	fmt.Println("string of dst bytes", string(dst.Bytes()))

	if err := json.Unmarshal(dst.Bytes(), &token); err != nil {
		log.Println("unmarshal error", err)
		return err
	}
	c.token.tokenString = token
	return nil
}

// GetProcessID gets processID as per provided search criteria
func (c *Client) GetProcess(ctx context.Context, searchCriteria domain.Criteria) (*domain.Process, error) {
	log.Println("GetProcess hit")
	processesURL := c.getProcessesURL
	if processesURL == "" {
		return nil, domain.ErrProcessesURLNotFound
	}
	url := processesURL + searchCriteria.ID
	req, err := http.NewRequest(http.MethodGet, url, nil)
	//req, err := http.NewRequest(http.MethodGet, "https://halykbpm-api.halykbank.nb/process-searcher/instance?searchValue=980124450084", nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	res, err := c.getDataWithRetries(req)
	if err != nil {
		return nil, err
	}
	var dst bytes.Buffer
	if _, err = io.Copy(&dst, res.Body); err != nil {
		log.Println("Error here")
		return nil, err
	}

	defer res.Body.Close()
	log.Println("Printing code sttus", res.StatusCode)
	var processes []*domain.Process
	if err := json.Unmarshal(dst.Bytes(), &processes); err != nil {
		//fmt.Println("Unramshal json error", err)
		return nil, err
	}
	if len(processes) == 0 {
		return nil, domain.ErrNoDataFound
	}

	// iterate process list to get one we need
	for _, process := range processes {
		if process.ProcessDefinitionKey == searchCriteria.Type {
			return process, nil
		}
	}
	return nil, domain.ErrProcessNotFound
}

// GetProcessStatus gets status of given process
func (c *Client) GetProcessStatus(ctx context.Context, processID string) (string, error) {
	log.Println("GetProcessStatus hit")
	if processID == "" {
		return "", domain.ErrProcessIDNotFound
	}
	processURL := c.getProcessURL
	if processURL == "" {
		return "", domain.ErrProcessURLNotFound
	}

	var jsonStr = []byte(fmt.Sprintf("{\"processInstanceIdIn\":[\"%s\"]}", processID))
	dst := bytes.NewBuffer(jsonStr)
	req, err := http.NewRequest(http.MethodPost, processURL, dst)
	req = req.WithContext(ctx)
	res, err := c.getDataWithRetries(req)
	if err != nil {
		return "", err
	}

	if _, err = io.Copy(dst, res.Body); err != nil {
		return "", err
	}

	defer res.Body.Close()

	var processVars []domain.ProcessStatus
	if err := json.Unmarshal(dst.Bytes(), &processVars); err != nil {
		return "", err
	}
	if len(processVars) < 40 {
		return "", domain.ErrNoVarsFound
	}
	for i := 35; i <= 40; i++ {
		message := processVars[i]
		if value, ok := message.Value.(string); ok {
			if message.Name == "processStatusMessage" { // index 195!enkp
				if value != "" {
					return value, nil
				}
			}
			if len(value) > 4 {
				fmt.Println(message)
				return value, nil
			}
		}
	}

	return "", domain.ErrProcessStatusNotFound
}

func (c *Client) getToken() string {
	c.token.mx.RLock()
	defer c.token.mx.RUnlock()
	token := c.token.tokenString
	return token
}

// getDataWithRetries attempts to retrieve data based on given request and returns response, if any
func (c *Client) getDataWithRetries(req *http.Request) (*http.Response, error) {
	log.Println("getDatawithRetries hit")
	var res *http.Response
	var err error
	var status int
	token := c.getToken()
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	for i := 0; i < 5; i++ {
		fmt.Println("Retry#", i)
		if token == "" || status == 401 {
			log.Println("reset token")
			if err := c.setToken(); err != nil {
				log.Println("401")
				return nil, err
			}
			token = c.getToken()
			req.Header.Set("Authorization", "Bearer "+token)
		}

		res, err = c.Do(req)
		if err != nil {
			log.Println("Do req err")
			return nil, err
		}
		status = res.StatusCode
		if res.StatusCode != 401 {
			break
		}
	}
	if status == 401 {
		return nil, domain.ErrUnauthorized
	}
	return res, nil
}

// RetryJobOrTask retries job or external task when incidents occur
func (c *Client) RetryJobOrTask(ctx context.Context, processID string) error {
	log.Println("Retry Job or task")
	if processID == "" {
		return domain.ErrProcessIDNotFound
	}
	incidentURL := c.getIncidentsURL
	if incidentURL == "" {
		return domain.ErrIncidentURLNotFound
	}
	incidentURL += processID
	fmt.Println("incident url", incidentURL)
	req, err := http.NewRequest(http.MethodGet, incidentURL, nil)
	req = req.WithContext(ctx)

	res, err := c.getDataWithRetries(req)
	if err != nil {
		return err
	}
	var dst bytes.Buffer
	if _, err = io.Copy(&dst, res.Body); err != nil {
		return err
	}

	res.Body.Close()

	var incidents []domain.Incident
	if err := json.Unmarshal(dst.Bytes(), &incidents); err != nil {
		return err
	}
	if len(incidents) == 0 {
		fmt.Println("len0")
		return domain.ErrNoIncidentFound
	}
	log.Printf("PRINTIN INCIDENT%#v", incidents[0])
	log.Println("LENGTH", len(incidents))
	incident := incidents[0]
	incidentType := incident.IncidentType
	log.Println("Configuration:", incident.Configuration)
	var retriesURL string
	jsonStr := []byte("{\"retries\":1}")
	log.Println("incident type:", incidentType)
	if incidentType == "failedExternalTask" {
		retriesURL = c.retryTaskURL
		if retriesURL == "" {
			return domain.ErrRetriesURLNotFound
		}
		jsonStr = []byte(fmt.Sprintf("{\"retries\":1, \"externalTaskIds\":[\"%s\"]}", incident.Configuration))
	} else if incidentType == "failedJob" {
		retriesURL = c.retryJobURL
		if retriesURL == "" {
			return domain.ErrRetriesURLNotFound
		}
		retriesURL += fmt.Sprintf("/%s/retries", incident.Configuration)
	} else {
		return domain.ErrUnknownIncident
	}
	fmt.Println(string(jsonStr), "\n", retriesURL)
	req, err = http.NewRequest(http.MethodPut, retriesURL, bytes.NewBuffer(jsonStr))
	req = req.WithContext(ctx)

	res, err = c.getDataWithRetries(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	// TODO: проверить есть ли все еще инциденты и несколько раз ретрай если есть
	log.Println("Retry status:", res.StatusCode)
	return nil
}

// GetActivityID gets ID for a given activity name
func (c *Client) GetActivityID(ctx context.Context, processID, activityName string) (string, error) {
	if processID == "" {
		return "", domain.ErrProcessIDNotFound
	}
	url := c.activitySearchURL
	if url == "" {
		return "", domain.ErrActivitySearchURLNotFound
	}
	url += processID
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req = req.WithContext(ctx)
	res, err := c.getDataWithRetries(req)
	if err != nil {
		return "", err
	}
	var dst bytes.Buffer
	if _, err = io.Copy(&dst, res.Body); err != nil {
		log.Println("Error here")
		return "", err
	}

	defer res.Body.Close()
	log.Println("Printing code sttus", res.StatusCode)
	var activities []domain.Activity
	if err := json.Unmarshal(dst.Bytes(), &activities); err != nil {
		//fmt.Println("Unramshal json error", err)
		return "", err
	}
	if len(activities) == 0 {
		return "", domain.ErrNoActivityFound
	}
	var ID string
	for _, a := range activities {
		if a.ActivityName == activityName {
			ID = a.Id
			if ID == "" {
				return "", domain.ErrActivityIDNotFound
			}
			return ID, nil
		}
	}
	return "", domain.ErrActivityNotFound
}

// Redo reattempts an activity, e.g. UVK
func (c *Client) Redo(ctx context.Context, processID, activityID string) error {
	log.Println("Redo")
	if processID == "" {
		return domain.ErrProcessIDNotFound
	}
	if activityID == "" {
		return domain.ErrActivityIDNotFound
	}
	url := c.modificationURL
	if url == "" {
		return domain.ModificationURLNotFound
	}
	url += fmt.Sprintf("%s/modification", processID)
	// Request body

	//{"instructions":[{"transitionId":"Flow_189texq","type":"startTransition"},{"type":"cancel","activityInstanceId":"notificationEnd:gf3e1b6e0f94d022"}],"skipCustomListeners":true,"skipIoMappings":true}
	var jsonStr = []byte(fmt.Sprintf("{\"instructions\":[{\"transitionId\":\"Flow_189texq\", \"type\":\"startTransition\"},{\"type\":\"cancel\", \"activityInstanceId\":\"%s\"}], \"skipCustomListeners\":true,\"skipIoMappings\":true}", processID))
	dst := bytes.NewBuffer(jsonStr)
	req, err := http.NewRequest(http.MethodPost, url, dst)
	fmt.Println("modification url", url)
	req = req.WithContext(ctx)

	res, err := c.getDataWithRetries(req)
	if err != nil {
		return err
	}
	// defer res.Body.Close() // response isn't supposed to contain body
	log.Println("Modification status:", res.StatusCode)
	return nil
}

// UpdateBranch updates two variables: branchSapCode and initRole
func (c *Client) UpdateBranch(ctx context.Context, processID, branchCode string) error {
	url := c.updateVarsURL
	if url == "" {
		return domain.ErrUpdateVarsURLNotFound
	}
	if processID == "" {
		return domain.ErrProcessIDNotFound
	}
	url += fmt.Sprintf("%s/localVariables", processID)
	var jsonStr = []byte(fmt.Sprintf("{\"modifications\": {\"branchSapCode\": {\"type\": \"String\", \"value\": \"%s\"}}},  {\"initRole\": {\"type\": \"String\", \"value\": \"app-front-cashier-vk-%s\"}}}", branchCode, branchCode))
	dst := bytes.NewBuffer(jsonStr)
	req, err := http.NewRequest(http.MethodPost, url, dst)
	fmt.Println("modification url", url)

	req = req.WithContext(ctx)

	res, err := c.getDataWithRetries(req)
	if err != nil {
		return err
	}
	fmt.Println("update branchSapCode status", res.Status)
	if res.StatusCode != 204 {
		return domain.ErrUpdateFailed
	}
	return nil
	//jsontStr =
}

// GetRole fetches role of manager sending request
func (c *Client) GetRole(ctx context.Context, tab string) error {
	if tab == "" {
		return domain.ErrInvalidTab
	}
	url := c.managerRoleURL
	if url == "" {
		return domain.ErrRoleURLNotFound
	}
	url += tab
	req, err := http.NewRequest(http.MethodGet, url, nil)
	//req, err := http.NewRequest(http.MethodGet, "https://halykbpm-api.halykbank.nb/process-searcher/instance?searchValue=980124450084", nil)
	if err != nil {
		return err
	}
	req = req.WithContext(ctx)
	res, err := c.getDataWithRetries(req)
	if err != nil {
		return err
	}
	var dst bytes.Buffer
	if _, err = io.Copy(&dst, res.Body); err != nil {
		log.Println("Error here")
		return err
	}

	defer res.Body.Close()
	log.Println("Printing code sttus", res.StatusCode)
	var roles []domain.Role
	if err := json.Unmarshal(dst.Bytes(), &roles); err != nil {
		//fmt.Println("Unramshal json error", err)
		return err
	}
	if len(roles) == 0 {
		return domain.ErrNoRoleFound
	}
	return nil
}
