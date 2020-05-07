package crocgodyl

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// VERSION of crocgodyl follows Semantic Versioning. (http://semver.org/)
const VERSION = "0.0.4-alpha"

// ErrorResponse is the response from the panels for errors.
type ErrorResponse struct {
	Errors []Error `json:"errors"`
}

// Error is the struct from the panels for errors.
type Error struct {
	Code   string `json:"code,omitempty"`
	Status string `json:"status,omitempty"`
	Detail string `json:"detail,omitempty"`
	Source struct {
		Field string `json:"field,omitempty"`
	} `json:"source,omitempty"`
}

// Meta is the meta information on some queries
type Meta struct {
	Pagination Pagination `json:"pagination,omitempty"`
}

// Pagination is the information how many responses there on a page and how many pages there are.
type Pagination struct {
	Total       int   `json:"total,omitempty"`
	Count       int   `json:"count,omitempty"`
	PerPage     int   `json:"per_page,omitempty"`
	CurrentPage int   `json:"current_page,omitempty"`
	TotalPages  int   `json:"total_pages,omitempty"`
	Links       Links `json:"links"`
}

// Links is the struct for the links in the Pagination struct
type Links struct {
	Previous string `json:"previous,omitempty"`
	Next     string `json:"next,omitempty"`
}

//
// crocgodyl
//

// CrocConfig is the config for crocgodyl
type CrocConfig struct {
	PanelURL    string
	ClientToken string
	APIToken    string
}

//
// Application code
//

// NewCrocConfig sets up the API interface with
func NewCrocConfig(panelURL string, clientToken string, appToken string) (config *CrocConfig, err error) {

	if panelURL == "" && clientToken == "" && appToken == "" {
		return config, errors.New("you need to configure the panel and at least one api token")
	}

	if panelURL == "" {
		return config, errors.New("a panel URL is required to use the API")
	}

	if clientToken == "" && appToken == "" {
		return config, errors.New("at least one api token is required")
	}

	config = &CrocConfig{}
	config.PanelURL = panelURL
	config.ClientToken = clientToken
	config.APIToken = appToken

	// validate the server is up and available
	if _, err = config.GetUsers(); err != nil {
		return config, nil
	}

	return
}

func (config *CrocConfig) queryPanelAPI(endpoint, request string, data []byte) ([]byte, error) {
	return config.queryPanelCallback("application", config.APIToken, endpoint, request, data)
}

func (config *CrocConfig) queryPanelClient(endpoint, request string, data []byte) ([]byte, error) {
	return config.queryPanelCallback("client", config.ClientToken, endpoint, request, data)
}

func (config *CrocConfig) queryPanelCallback(sector, token, endpoint, request string, data []byte) ([]byte, error) {
	var bodyBytes []byte

	client := &http.Client{}
	url := fmt.Sprintf("%s/api/%s/%s", config.PanelURL, sector, endpoint)

	var req *http.Request
	switch strings.ToLower(request){
	case "get":
		req, _  = http.NewRequest("GET", url, nil)
	case "post":
		req, _ = http.NewRequest("POST", url, bytes.NewBuffer(data))
	case "patch":
		req, _ = http.NewRequest("PATCH", url, bytes.NewBuffer(data))
	case "delete":
		req, _ = http.NewRequest("DELETE", url, nil)
	default:
		return nil, errors.New("method not allowed")
	}

	//Sets request header for the http request
	req.Header.Add("Authorization", "Bearer " + token)
	req.Header.Add("Accept", "Application/vnd.pterodactyl.v1+json")
	req.Header.Set("Content-Type", "application/json")

	//send request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode > 200 && resp.StatusCode < 204 {
		bodyBytes, _ = ioutil.ReadAll(resp.Body)
		return bodyBytes, errors.New(string(bodyBytes))
	}

	if resp.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(resp.Body)
	}

	//set bodyBytes to the response body
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	//Close response thread
	err =  resp.Body.Close()
	if err != nil {
		return bodyBytes, errors.New("unable to close response body")
	}

	//return byte structure
	return bodyBytes, nil
}

func printJSON(b []byte) ([]byte, error) {
	var out bytes.Buffer
	err := json.Indent(&out, b, "", "  ")
	return out.Bytes(), err
}
