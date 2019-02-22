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
	req.Header.Set("Content-Type", "application/json")

	//send request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 403 {
		bodyBytes, _ = ioutil.ReadAll(resp.Body)
		return bodyBytes, errors.New("the request was valid, but the server is refusing action")
	}

	if resp.StatusCode == 404 {
		bodyBytes, _ = ioutil.ReadAll(resp.Body)
		return bodyBytes, errors.New("the requested resource could not be found but may be available in the future")
	}

	if resp.StatusCode == 422 {
		bodyBytes, _ = ioutil.ReadAll(resp.Body)
		return bodyBytes, errors.New("the request was well-formed but was unable to be followed due to semantic errors")
	}

	if resp.StatusCode >= 500 {
		bodyBytes, _ = ioutil.ReadAll(resp.Body)
		return bodyBytes, errors.New("the server failed to fulfil a request")
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
