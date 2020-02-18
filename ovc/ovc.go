// Package ovc implements a client library for SimpliVity REST API endpoints.
package ovc

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// Logs will be helpful for debugging purpose
// Comment the init function to get the logs
func init() {
	log.SetOutput(ioutil.Discard)
}

// Http status codes of the OVC API endpoints.
var (
	httpStatusCodes = map[int]bool{
		http.StatusOK:                    true,
		http.StatusCreated:               true,
		http.StatusAccepted:              true,
		http.StatusNoContent:             true,
		http.StatusBadRequest:            false,
		http.StatusUnauthorized:          false,
		http.StatusNotFound:              false,
		http.StatusMethodNotAllowed:      false,
		http.StatusRequestEntityTooLarge: false,
		http.StatusUnsupportedMediaType:  false,
		http.StatusInternalServerError:   false,
		http.StatusBadGateway:            false,
		http.StatusGatewayTimeout:        false,
	}

	commonClient *Client
)

// Query Parameters of the Get all endpoints.
type GetAllParams struct {
	Limit              int
	Offset             int
	Sort               string
	Order              string
	Fields             string
	Case               string
	ShowOptionalFields bool
	Filters            map[string]string
}

// QueryString creates query string from the GetAllParams parameters including filters
func (p GetAllParams) QueryString() string {
	QueryStr := url.Values{}

	if p.Limit < 1 {
		p.Limit = 500
	}
	QueryStr.Add("limit", strconv.Itoa(p.Limit))
	QueryStr.Add("offset", strconv.Itoa(p.Offset))
	QueryStr.Add("sort", p.Sort)
	QueryStr.Add("order", p.Order)
	QueryStr.Add("fields", p.Fields)
	QueryStr.Add("case", p.Case)
	QueryStr.Add("show_optional_fields", strconv.FormatBool(p.ShowOptionalFields))

	if len(p.Filters) > 0 {
		for key, value := range p.Filters {
			QueryStr.Add(key, value)
		}
	}
	return QueryStr.Encode()
}

// Client handles communications with the OVC API.
//
// SimpliVity API doc: https://developer.hpe.com/api/simplivity/
type Client struct {
	client *http.Client
	common resourceClient // Common OVC client for all the resources

	// OVC IP
	OVCIP string

	//OVC username
	Username string

	// OVC password
	Password string

	// OVC access token
	AccessToken string

	// SSL certificate path
	SSLCertificatePath string

	//OVC resource clients
	Backups           *BackupResource
	Datastores        *DatastoreResource
	Hosts             *HostResource
	OmniStackClusters *OmniStackClusterResource
	Policies          *PolicyResource
	VirtualMachines   *VirtualMachineResource
	Tasks             *TaskResource
}

// Helps to share a common OVC client with all the resource clients.
type resourceClient struct {
	client *Client
}

// Auth token endpoint response.
type auth struct {
	AccessToken string `json:"access_token,omitempty"`
}

// NewClient creates a new OVC client.
// Sets http client to make http connection with the OVC.
// Sets access token by communicating with the auth token endpoint.
// Initializes resource clients with a common OVC client.
func NewClient(username string, password string, ovc_ip string, ssl_certificate string) (*Client, error) {
	c := &Client{Username: username, Password: password, OVCIP: ovc_ip, SSLCertificatePath: ssl_certificate}
	c.common.client = c

	// Set Http client.
	err := c.setHttpClient()
	if err != nil {
		return nil, err
	}

	commonClient = c

	// Login and get access token using the username and password.
	err = c.SetAccessToken()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// Initialize resource clients.
	c.Backups = (*BackupResource)(&c.common)
	c.Datastores = (*DatastoreResource)(&c.common)
	c.Hosts = (*HostResource)(&c.common)
	c.OmniStackClusters = (*OmniStackClusterResource)(&c.common)
	c.VirtualMachines = (*VirtualMachineResource)(&c.common)
	c.Policies = (*PolicyResource)(&c.common)
	c.Tasks = (*TaskResource)(&c.common)

	return c, nil
}

// setHttpClient creates a http client to communicate with the OVC.
// A secured connection will be established, if a SSL certificate file path exists.
func (c *Client) setHttpClient() error {

	var tr *http.Transport

	if c.SSLCertificatePath != "" {
		caCert, err := ioutil.ReadFile(c.SSLCertificatePath)
		if err != nil {
			//log.Fatal(err)
			return err
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		tr = &http.Transport{
			TLSClientConfig: &tls.Config{RootCAs: caCertPool}}

	} else {
		tr = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	}

	//Set Http client to make requests to OVC
	c.client = &http.Client{Transport: tr}

	return nil
}

// CreateResourceURL creats the resource URL by appending the base URL with the resource path
func (c *Client) CreateResourceURL(path string, query string) (*url.URL, error) {

	baseURL := c.OVCIP
	if !strings.HasPrefix(c.OVCIP, "http") {
		baseURL = fmt.Sprintf("https://%s/api", c.OVCIP)
	}

	resourceURL, err := url.Parse(baseURL + path)
	if err != nil {
		//log.Fatal(err)
		return nil, err
	}

	if query != "" {
		resourceURL.RawQuery = query
	}

	return resourceURL, nil
}

// SetAccessToken sets the access token which will be used by other endpoints
// Makes a call to the login endpoint and gets the token.
func (c *Client) SetAccessToken() error {
	endpoint, err := c.CreateResourceURL("/oauth/token", "")
	if err != nil {
		log.Println(err)
		return err
	}

	//Sets form data of the login API
	reqData := url.Values{}
	reqData.Set("username", c.Username)
	reqData.Set("password", c.Password)
	reqData.Set("grant_type", "password")

	//Creates new http request with URL encoded body
	req, err := http.NewRequest("POST", endpoint.String(), bytes.NewBufferString(reqData.Encode()))
	if err != nil {
		return err
	}

	//Login API expects URL encoded body
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")

	//Basic auth should be set with username simplivity and empty password
	req.SetBasicAuth("simplivity", "")

	//Makes request to the login API
	data, err, httpError := c.Do(req)
	if err != nil {
		return err
	}

	// Handles http errors.
	if httpError != nil {
		return fmt.Errorf("Error: Status code:%s - %s", httpError.Status, httpError.Message)
	}

	//Unmarshal the login API response to get the token
	authData := auth{}
	err = json.Unmarshal(data, &authData)
	if err != nil {
		return err
	}

	//Set access token
	c.AccessToken = authData.AccessToken

	return nil
}

// SetHeaders sets headers of the API requests.
func (c *Client) SetHeaders(req *http.Request, headers map[string]string) {
	//Required for all the API requests except login API
	req.Header.Set("Authorization", "Bearer "+c.AccessToken)

	//Set other headers of the API request
	for name, value := range headers {
		req.Header.Set(name, value)
	}
}

// DoRequest creates a new http request and make calls to the OVC.
// Creates a new http request using NewRequest method.
// Makes http call to the OVC using Do method.
// Tries to get another token and makes a fresh request if the token expired.
func (c *Client) DoRequest(method, path, queryStr string, body interface{}, headers map[string]string) ([]byte, error) {
	var data []byte

	req, err := c.NewRequest(method, path, queryStr, body)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	req_headers := map[string]string{}
	// Setting default headers for API request
	if body != nil {
		req_headers["Content-Type"] = "application/vnd.simplivity.v1+json"
	}

	if headers != nil {
		for name, value := range headers {
			req_headers[name] = value
		}
	}

	c.SetHeaders(req, req_headers)
	data, err, httpError := c.Do(req)

	// Get a fresh token and make a new request if the current token is expired.
	if httpError != nil {
		if httpError.Error == "invalid_token" {
			log.Println("Token expired - trying to make another call with a new token")
			c.SetAccessToken()
			c.SetHeaders(req, req_headers)
			data, err, httpError = c.Do(req)
			if httpError != nil {
				err = fmt.Errorf("Error: Status code:%s - %s", httpError.Status, httpError.Message)
			}
			return data, err
		}
		err = fmt.Errorf("Error: Status code:%s - %s", httpError.Status, httpError.Message)
	}

	return data, err
}

// NewRequest creates a new http request by setting the URL, method and body.
func (c *Client) NewRequest(method, path, query string, body interface{}) (*http.Request, error) {
	resourceURL, err := c.CreateResourceURL(path, query)
	if err != nil {
		log.Println(err)
	}

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, resourceURL.String(), buf)
	if err != nil {
		return nil, err
	}

	return req, nil
}

// Error response from an API endpoint.
type OVCRespError struct {
	Exception string `json:"exception,omitempty"`
	Path      string `json:"path,omitempty"`
	Error     string `json:"error,omitempty"`
	Message   string `json:"message,omitempty"`
	Status    string `json:"status,omitempty"`
}

// Do makes calls to the OVC
// Returns data if the call is successfull or error
func (c *Client) Do(req *http.Request) ([]byte, error, *OVCRespError) {
	resp, err := c.client.Do(req)

	if err != nil {
		return nil, err, nil
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err, nil
	}

	if !httpStatusCodes[resp.StatusCode] {
		log.Println("HTTP status code", resp.StatusCode)
		var errResp OVCRespError

		err = json.Unmarshal(data, &errResp)
		if err != nil {
			return nil, err, nil
			log.Println("Unmarshal error", err)
		}

		return nil, nil, &errResp
	}

	return data, nil, nil
}
