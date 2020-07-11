package crocgodyl

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
)

func (config *AppConfig) queryApplicationAPI(endpoint, request string, data []byte) (bodyBytes []byte, err error) {
	//http get json request
	client := &http.Client{}
	req, _ := http.NewRequest("GET", config.PanelURL+"/api/application/"+endpoint, nil)

	switch {
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
		return
	}

	if resp.StatusCode != 200 && resp.StatusCode != 201 && resp.StatusCode != 202 && resp.StatusCode != 204 {
		bodyBytes, _ = ioutil.ReadAll(resp.Body)
		err = errors.New(string(bodyBytes))
		return
	}

	if resp.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(resp.Body)
	}

	//Close response thread
	defer resp.Body.Close()

	bodyBytes = bytes.Replace(bodyBytes, []byte("[]"), []byte("{}"), -1)

	//return byte structure
	return
}