package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"support/domain"
	"sync"
	"time"

	"github.com/Azure/go-ntlmssp"
)

type Client struct {
	clientOnce sync.Once
	builder    *clientBuilder
	http.Client
	token    string
	username string
	password string
	http.Request
	ntlmssp.Negotiator
}

type clientBuilder struct {
	maxIdleConnections int
	connectionTimeout  time.Duration
	responseTimeout    time.Duration
	disableTimeouts    bool
	client             *http.Client
	userAgent          string
}

const (
	defaultMaxIdleConnections = 5
	defaultResponseTimeout    = 5 * time.Second
	defaultConnectionTimeout  = 1 * time.Second
)

func NewClient() *Client {
	c := &Client{
		username: os.Getenv("USERNAME"),
		password: os.Getenv("PASSWORD"),
	}
	c.Transport = ntlmssp.Negotiator{
		RoundTripper: &http.Transport{
			MaxIdleConnsPerHost: 5,
			DialContext: (&net.Dialer{
				Timeout: time.Second * 5,
			}).DialContext,
		},
	}
	return c

}
func NewBuilder() ClientBuilder {
	builder := &clientBuilder{}
	return builder
}

type ClientBuilder interface {
	SetConnectionTimeout(timeout time.Duration) ClientBuilder
	SetResponseTimeout(timeout time.Duration) ClientBuilder
	SetMaxIdleConnections(i int) ClientBuilder
	DisableTimeouts(disable bool) ClientBuilder
	// SetHttpClient(c *http.Client) ClientBuilder
	SetUserAgent(userAgent string) ClientBuilder
	Build() domain.ProcessRepository
}

func (c *clientBuilder) Build() domain.ProcessRepository {
	return &Client{
		builder: c,
	}
}

func (c *clientBuilder) SetConnectionTimeout(timeout time.Duration) ClientBuilder {
	c.connectionTimeout = timeout
	return c
}

func (c *clientBuilder) SetResponseTimeout(timeout time.Duration) ClientBuilder {
	c.responseTimeout = timeout
	return c
}

func (c *clientBuilder) SetMaxIdleConnections(i int) ClientBuilder {
	c.maxIdleConnections = i
	return c
}

func (c *clientBuilder) DisableTimeouts(disable bool) ClientBuilder {
	c.disableTimeouts = disable
	return c
}

func (c *clientBuilder) SetHttpClient(client *http.Client) ClientBuilder {
	c.client = client
	return c
}

func (c *clientBuilder) SetUserAgent(userAgent string) ClientBuilder {
	c.userAgent = userAgent
	return c
}

// // NewClient returns new Client
// func NewClient() *Client {
// 	c := &Client{
// 		username: os.Getenv("USERNAME"),
// 		password: os.Getenv("PASSWORD"),
// 	}
// 	//c.Transport =
// 	//Transport: &http.Transport {MaxIdleConnsPerHost: :5}

// }

// SetToken allows user set their own token
func (c *Client) SetToken(token string) error {
	if token == "" {
		return domain.ErrTokenNotFound
	}
	c.token = token
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

	fmt.Println("string of dst bytes", string(dst.Bytes()))
	if err := json.Unmarshal(dst.Bytes(), c.token); err != nil {
		log.Println("unmarshal error", err)
		return err
	}
	return nil
}

// getProcessID gets processID as per provided search criteria
func (c *Client) GetProcess(ctx context.Context, searchCriteria domain.Criteria) (*domain.Process, error) {
	processesURL := os.Getenv("PROCESSES_URL")
	if processesURL == "" {
		return nil, domain.ErrProcessesURLNotFound
	}
	req, err := http.NewRequest(http.MethodGet, processesURL+searchCriteria.ID, nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	req.Header.Set("Authorization", "Bearer "+c.token)
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
		return nil, err
	}

	defer res.Body.Close()

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
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.token)

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
	for i, hing := range processVars {
		fmt.Println(i, hing)
	}

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

// getDataWithRetries attempts to retrieve data based on given request and returns response, if any
func (c *Client) getDataWithRetries(req *http.Request) (*http.Response, error) {
	var res *http.Response
	var err error
	var status int
	for i := 0; i < 5; i++ {
		if c.token == "" || status == 401 {
			if err := c.setToken(); err != nil {
				return nil, err
			}
		}

		res, err = c.Do(req)
		if err != nil {
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
	if processID == "" {
		return domain.ErrProcessIDNotFound
	}
	incidentURL := os.Getenv("GET_INCIDENT_URL")
	if incidentURL == "" {
		return domain.ErrIncidentURLNotFound
	}

	req, err := http.NewRequest(http.MethodGet, incidentURL, nil)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.token)

	res, err := c.getDataWithRetries(req)
	if err != nil {
		return err
	}
	var dst *bytes.Buffer
	if _, err = io.Copy(dst, res.Body); err != nil {
		return err
	}

	defer res.Body.Close()

	var incidents []domain.Incident
	//var A []Sth2
	if err := json.Unmarshal(dst.Bytes(), &incidents); err != nil {
		//fmt.Println("err last req unmarshal", err)
		return err
	}
	if len(incidents) == 0 {
		return domain.ErrNoIncidentFound
	}
	incident := incidents[0]
	incidentType := incident.IncidentType
	var retriesURL string
	var jsonStr []byte
	if incidentType == "failedExternalTask" {
		retriesURL = os.Getenv("RETRY_TASK_URL")
		if retriesURL == "" {
			return domain.ErrRetriesURLNotFound
		}
		jsonStr := []byte(fmt.Sprintf("{\"retries\":1, \"externalTaskIds\":[\"%s\"]}", incident.Configuraion))
		dst := bytes.NewBuffer(jsonStr)
		req, err = http.NewRequest(http.MethodPut, retriesURL, dst)
	} else if incidentType == "failedJob" {
		retriesURL = os.Getenv("RETRY_JOB_URL")
		if retriesURL == "" {
			return domain.ErrRetriesURLNotFound
		}
		retriesURL += fmt.Sprintf("/%s/retries", incident.Configuraion)
	} else {
		return domain.ErrUnknownIncident
	}

	req, err = http.NewRequest(http.MethodPut, retriesURL, nil)
	req = req.WithContext(ctx)
	if jsonStr != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	res, err = c.getDataWithRetries(req)
	if err != nil {
		return err
	}
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
