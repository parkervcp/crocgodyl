package crocgodyl

import (
	"errors"
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

type Links struct {
	Next     string `json:"next"`
	Previous string `json:"previous"`
}

//
// crocgodyl
//

// AppConfig is the config for crocgodyl
type AppConfig struct {
	PanelURL    string
	AppToken    string
}

type ClientConfig struct {
	PanelURL    string
	ClientToken string
}

//
// Application code
//

// NewCrocConfig sets up the API interface with
func NewAppClient(panelURL string, appToken string) (config *AppConfig, err error) {

	config = &AppConfig{}

	if panelURL == "" && appToken == "" {
		return config, errors.New("you need to configure the panel and at least one api token")
	}

	if panelURL == "" {
		return config, errors.New("a panel URL is required to use the API")
	}

	if appToken == "" {
		return config, errors.New("an application token is required")
	}

	config.PanelURL = panelURL
	config.AppToken = appToken

	// validate the server is up and available
	if _, err = config.getUserByPage(1); err != nil {
		return config, err
	}

	return
}

// NewCrocConfig sets up the API interface with
func NewClientClient(panelURL string, clientToken string) (config *ClientConfig, err error) {

	config = &ClientConfig{}

	if panelURL == "" && clientToken == "" {
		return config, errors.New("you need to configure the panel and at least one api token")
	}

	if panelURL == "" {
		return config, errors.New("a panel URL is required to use the API")
	}

	if clientToken == "" {
		return config, errors.New("an application token is required")
	}

	config.PanelURL = panelURL
	config.ClientToken = clientToken

	// validate the server is up and available
	if _, err = config.getClientServersByPage(1); err != nil {
		return config, err
	}

	return
}