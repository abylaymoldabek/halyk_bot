package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"v/domain"
	"v/ntlmssp"
	"net"
	"time"
	"sync"
)

type Token struct {
	tokenString string
	mx sync.RWMutex
} 

type Client struct {
	http.Client
	token   Token
	username string
	password string
}

// NewClient returns new Client
func NewClient() *Client {
	c := &Client{ 
		username: os.Getenv("USERNAME"),
		password: os.Getenv("PASSWORD"),
	}
	c.setToken()
	c.Transport = ntlmssp.Negotiator{
		RoundTripper: &http.Transport{
			MaxIdleConnsPerHost: 5,
			DialContext:(&net.Dialer{
				Timeout: time.Second*5,
			}).DialContext,
		},
	}
	return c
	
}
// maxIdleConnections int
// 	connectionTimeout  time.Duration
// 	responseTimeout    time.Duration
// 	disableTimeouts    bool
// 	baseUrl            string
// 	client             *http.Client
// 	userAgent          string

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
	tokenURL := os.Getenv("TOKEN_URL")
	if tokenURL == "" {
		return domain.ErrTokenURLNotFound
	}
	
	req, err := http.NewRequest(http.MethodGet, tokenURL, nil)
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

// getProcessID gets processID as per provided search criteria
func (c *Client) GetProcess(ctx context.Context, searchCriteria domain.Criteria) (*domain.Process, error) {
	log.Println("GetProcess hit")
	processesURL := os.Getenv("PROCESSES_URL")
	if processesURL == "" {
		return nil, domain.ErrProcessesURLNotFound
	}
	url := processesURL+searchCriteria.ID
	req, err := http.NewRequest(http.MethodGet, url, nil)
	//req, err := http.NewRequest(http.MethodGet, "https://halykbpm-api.halykbank.nb/process-searcher/instance?searchValue=980124450084", nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	// c.token.mx.RWLock()
	// defer c.token.mx.RWUnlock()
	// req.Header.Set("Authorization", "Bearer "+c.token)
	// TODO: IF EVERYTHING OK, DELETE 3 LINES BELOW!!!!
	// res, err := c.Do(req)
	// if err != nil {
	// 	return nil, err
	// }

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
	log.Println("Printing code sttus",res.StatusCode)
	//fmt.Println("dst bytes:\n\n\n\n\n", string(dst.Bytes()))
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


// getProcessStatus gets status of given process
func (c *Client) GetProcessStatus(ctx context.Context, processID string) (string, error) {
	log.Println("GetProcessStatus hit")
	if processID == "" {
		return "", domain.ErrProcessIDNotFound
	}
	processURL := os.Getenv("PROCESS_URL")
	if processURL == "" {
		return "", domain.ErrProcessURLNotFound
	}

	var jsonStr = []byte(fmt.Sprintf("{\"processInstanceIdIn\":[\"%s\"]}", processID))
	dst := bytes.NewBuffer(jsonStr)
	// var jsonStr = []byte(fmt.Sprintf("{\"processInstanceIdIn\":[\"%s\"]}", processID))
	req, err := http.NewRequest(http.MethodPost, processURL, dst)
	req = req.WithContext(ctx)
	//c.token.mx.RLock()
	// defer c.token.mx.RWUnlock()
	// req.Header.Set("Content-Type", "application/json")
	// req.Header.Set("Authorization", "Bearer "+c.token)

	res, err := c.getDataWithRetries(req)
	if err != nil {
		return "", err
	}

	if _, err = io.Copy(dst, res.Body); err != nil {
		return "", err
	}

	defer res.Body.Close()

	var processVars []domain.ProcessStatus
	//var A []Sth2
	if err := json.Unmarshal(dst.Bytes(), &processVars); err != nil {
		//fmt.Println("err last req unmarshal", err)
		return "", err
	}
	if len(processVars) < 40 {
		return "", domain.ErrNoVarsFound
	}
	// TODO: Remove below
	// for i, hing := range processVars {
	// 	fmt.Println(i, hing)
	// }

	for i := 35; i <= 40; i++ {
		message := processVars[i]
		if value, ok := message.Value.(string); ok && len(value) > 4 {
			//if value := message.Value; len(value)>4 && value[:4] == "done" {
			fmt.Println(message)
			return value, nil
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
		fmt.Println("Retry#",i)
		//c.token.mx.RLock()
		if token == "" || status == 401 {
		//	c.token.mx.RUnlock()
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
		if res.StatusCode != 401 { //TODO kakoi tam success code / unauthorized code???
			break
		}
	}
	if status == 401 { //TODO: or if status != 204 return nil, someerror
		return nil, domain.ErrUnauthorized
	}
	return res, nil
}

func (c *Client) RetryJobOrTask(ctx context.Context, processID string) error {
	log.Println("Retry Job or task")
	if processID == "" {
		return domain.ErrProcessIDNotFound
	}
	incidentURL := os.Getenv("GET_INCIDENT_URL")
	if incidentURL == "" {
		return domain.ErrIncidentURLNotFound
	}
	incidentURL+=processID
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
	//var A []Sth2
	if err := json.Unmarshal(dst.Bytes(), &incidents); err != nil {
		//fmt.Println("err last req unmarshal", err)
		return err
	}
	if len(incidents) == 0 {
		fmt.Println("len0")
		return domain.ErrNoIncidentFound
	}
	log.Printf("PRINTIN INCIDENT%#v", incidents[0])
	log.Println("LENGTH",len(incidents))
	incident := incidents[0]
	incidentType := incident.IncidentType
	log.Println("Configuration:", incident.Configuration)
	var retriesURL string
	jsonStr := []byte("{\"retries\":1}")
	log.Println("incident type:", incidentType)
	if incidentType == "failedExternalTask" {
		retriesURL = os.Getenv("RETRY_TASK_URL")
		if retriesURL == "" {
			return domain.ErrRetriesURLNotFound
		}
		jsonStr = []byte(fmt.Sprintf("{\"retries\":1, \"externalTaskIds\":[\"%s\"]}", incident.Configuration))
	} else if incidentType == "failedJob" {
		retriesURL = os.Getenv("RETRY_JOB_URL")
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

// func (c *Client) tokenValid() bool {
// 	tokenString := c.token
// 	if tokenString == "" {
// 		return false
// 	}
// 	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
// 	if err, ok := err.(*jwt.ValidationError); ok && err.Errors != jwt.ValidationErrorUnverifiable || !ok {
// 		fmt.Println(err)
// 		return false
// 	}

// 	if claims, ok := token.Claims.(jwt.MapClaims); ok {
// 		exp, ok := claims["exp"].(float64)
// 		if !ok {
// 			return false
// 		}
// 		expiredTime := time.Unix(int64(exp), 0)
// 		if time.Now().After(expiredTime) {
// 			return false
// 		}
// 	}
// 	return true
// }

// func (c *Client) getURLDataWithRetries(req *http.Request) (*http.Response, error) {
// 	var body []byte
// 	var err error
// 	var resp *http.Response

// 	for i := 0; i < 5; i++ {
// 		res, err := c.Do(req)

// 		if err == nil {
// 			break
// 		}

// 		fmt.Fprintf(os.Stderr, "Request error: %+v\n", err)
// 		fmt.Fprintf(os.Stderr, "Retrying in %v\n", backoff)
// 		time.Sleep(backoff)
// 	}

// 	// All retries failed
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	return resp, body, nil
// }
