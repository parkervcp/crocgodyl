package crocgodyl

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const Version = "0.0.3"

type Application struct {
	PanelURL string
	ApiKey   string
	Http     *http.Client
}

type Client struct {
	PanelURL string
	ApiKey   string
	Http     *http.Client
}

func NewApp(url, key string) (*Application, error) {
	if url == "" {
		return nil, errors.New("a valid panel url is required")
	}
	if key == "" {
		return nil, errors.New("a valid application api key is required")
	}

	app := Application{
		PanelURL: url,
		ApiKey:   key,
		Http:     &http.Client{},
	}

	return &app, nil
}

func (a *Application) newRequest(method, path string, body io.Reader) *http.Request {
	req, _ := http.NewRequest(method, fmt.Sprintf("%s/api/application%s", a.PanelURL, path), body)

	req.Header.Add("User-Agent", "Crocgodyl v"+Version)
	req.Header.Add("Authorization", "Bearer "+a.ApiKey)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	return req
}

func NewClient(url, key string) (*Client, error) {
	if url == "" {
		return nil, errors.New("a valid panel url is required")
	}
	if key == "" {
		return nil, errors.New("a valid client api key is required")
	}

	client := Client{
		PanelURL: url,
		ApiKey:   key,
		Http:     &http.Client{},
	}

	return &client, nil
}

func (a *Client) newRequest(method, path string, body io.Reader) *http.Request {
	req, _ := http.NewRequest(method, fmt.Sprintf("%s/api/client%s", a.PanelURL, path), body)

	req.Header.Add("User-Agent", "Crocgodyl v"+Version)
	req.Header.Add("Authorization", "Bearer "+a.ApiKey)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	return req
}

func validate(res *http.Response) ([]byte, error) {
	switch res.StatusCode {
	case http.StatusOK:
		fallthrough

	case http.StatusAccepted:
		fallthrough

	case http.StatusCreated:
		defer res.Body.Close()
		buf, _ := io.ReadAll(res.Body)
		return buf, nil

	case http.StatusNoContent:
		return nil, nil

	default:
		if res.StatusCode >= 500 {
			return nil, errors.New("internal server error")
		}

		defer res.Body.Close()
		buf, _ := io.ReadAll(res.Body)

		var errs *struct {
			Errors []ApiError
		}
		if err := json.Unmarshal(buf, errs); err != nil {
			return nil, err
		}

		return nil, &errs.Errors[0]
	}
}
