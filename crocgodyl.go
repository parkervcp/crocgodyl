package crocgodyl

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

// VERSION of crocgodyl follows Semantic Versioning. (http://semver.org/)
const VERSION = "0.0.2-alpha"

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
var config crocConfig

type crocConfig struct {
	PanelURL    string
	ClientToken string
	AppToken    string
}

//
// Application code
//

// New sets up the API interface with
func New(panelURL string, clientToken string, appToken string) error {

	if panelURL == "" && clientToken == "" && appToken == "" {
		return errors.New("You need to configure the panel and at least one api token")
	}

	if panelURL == "" {
		return errors.New("A panel URL is required to use the API")
	}

	if clientToken == "" && appToken == "" {
		return errors.New("At least one api token is required")
	}

	config.PanelURL = panelURL
	config.ClientToken = clientToken
	config.AppToken = appToken

	_, err := GetUsers()
	if err != nil {
		return err
	}

	return nil
}

func queryPanelAPI(endpoint string, request string, data []byte) ([]byte, error) {
	//var for response body byte
	var bodyBytes []byte
	//http get json request
	client := &http.Client{}
	req, _ := http.NewRequest("GET", config.PanelURL+"/api/application/"+endpoint, nil)

	switch {
	case request == "get":
	case request == "post":
		req, _ = http.NewRequest("POST", config.PanelURL+"/api/application/"+endpoint, bytes.NewBuffer(data))
	case request == "patch":
		req, _ = http.NewRequest("PATCH", config.PanelURL+"/api/application/"+endpoint, bytes.NewBuffer(data))
	case request == "delete":
		req, _ = http.NewRequest("DELETE", config.PanelURL+"/api/application/"+endpoint, nil)
	default:
	}

	//Sets request header for the http request
	req.Header.Add("Authorization", "Bearer "+config.AppToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	//send request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 && resp.StatusCode != 201 && resp.StatusCode != 202 && resp.StatusCode != 204 {
		bodyBytes, _ = ioutil.ReadAll(resp.Body)
		return bodyBytes, errors.New(string(bodyBytes))
	}

	if resp.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(resp.Body)
	}

	//set bodyBytes to the response body
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	//Close response thread
	defer resp.Body.Close()

	//return byte structure
	return bodyBytes, nil
}

func printJSON(b []byte) ([]byte, error) {
	var out bytes.Buffer
	err := json.Indent(&out, b, "", "  ")
	return out.Bytes(), err
}
