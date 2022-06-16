package crocgodyl

import (
	"bytes"
	"encoding/json"
	"time"
)

type Account struct {
	ID        int64  `json:"id"`
	Admin     bool   `json:"admin"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Language  string `json:"language"`
}

func (a *Account) FullName() string {
	return a.FirstName + " " + a.LastName
}

func (c *Client) Account() (*Account, error) {
	req := c.newRequest("GET", "/account", nil)
	res, err := c.Http.Do(req)
	if err != nil {
		return nil, err
	}

	buf, err := validate(res)
	if err != nil {
		return nil, err
	}

	var model *struct {
		Attributes Account `json:"attributes"`
	}
	if err = json.Unmarshal(buf, &model); err != nil {
		return nil, err
	}

	return &model.Attributes, nil
}

func (c *Client) TwoFactor() (string, error) {
	req := c.newRequest("GET", "/account/two-factor", nil)
	res, err := c.Http.Do(req)
	if err != nil {
		return "", err
	}

	buf, err := validate(res)
	if err != nil {
		return "", err
	}

	var model *struct {
		Data struct {
			URL string `json:"image_url_data"`
		} `json:"data"`
	}
	if err = json.Unmarshal(buf, &model); err != nil {
		return "", err
	}

	return model.Data.URL, nil
}

func (c *Client) EnableTwoFactor(code int) ([]string, error) {
	data, _ := json.Marshal(map[string]int{"code": code})
	body := bytes.Buffer{}
	body.Write(data)

	req := c.newRequest("POST", "/account/two-factor", &body)
	res, err := c.Http.Do(req)
	if err != nil {
		return nil, err
	}

	buf, err := validate(res)
	if err != nil {
		return nil, err
	}

	var model *struct {
		Attributes struct {
			Tokens []string `json:"tokens"`
		} `json:"attributes"`
	}
	if err = json.Unmarshal(buf, &model); err != nil {
		return nil, err
	}

	return model.Attributes.Tokens, nil
}

func (c *Client) DisableTwoFactor(password string) error {
	data, _ := json.Marshal(map[string]string{"password": password})
	body := bytes.Buffer{}
	body.Write(data)

	req := c.newRequest("DELETE", "/account/two-factor", &body)
	res, err := c.Http.Do(req)
	if err != nil {
		return err
	}

	if _, err = validate(res); err != nil {
		return err
	}

	return nil
}

func (c *Client) UpdateEmail(email, password string) error {
	data, _ := json.Marshal(map[string]string{"email": email, "password": password})
	body := bytes.Buffer{}
	body.Write(data)

	req := c.newRequest("PUT", "/account/email", &body)
	res, err := c.Http.Do(req)
	if err != nil {
		return err
	}

	if _, err = validate(res); err != nil {
		return err
	}

	return nil
}

func (c *Client) UpdatePassword(old, new string) error {
	data, _ := json.Marshal(map[string]string{
		"current_password":      old,
		"password":              new,
		"password_confirmation": new,
	})
	body := bytes.Buffer{}
	body.Write(data)

	req := c.newRequest("PUT", "/account/password", &body)
	res, err := c.Http.Do(req)
	if err != nil {
		return err
	}

	if _, err = validate(res); err != nil {
		return err
	}

	return nil
}

type ApiKey struct {
	Identifier  string    `json:"identifier"`
	Description string    `json:"description"`
	AllowedIPs  []string  `json:"allowed_ips"`
	CreatedAt   time.Time `json:"created_at"`
	LastUsedAt  time.Time `json:"last_used_at"`
}

func (c *Client) ApiKeys() ([]ApiKey, error) {
	req := c.newRequest("GET", "/account/api-keys", nil)
	res, err := c.Http.Do(req)
	if err != nil {
		return nil, err
	}

	buf, err := validate(res)
	if err != nil {
		return nil, err
	}

	var model *struct {
		Data []struct {
			Attributes ApiKey `json:"attributes"`
		} `json:"data"`
	}
	if err = json.Unmarshal(buf, &model); err != nil {
		return nil, err
	}

	keys := make([]ApiKey, 0, len(model.Data))
	for _, k := range model.Data {
		keys = append(keys, k.Attributes)
	}

	return keys, nil
}

func (c *Client) CreateKey(description string, ips []string) (*ApiKey, error) {
	data, _ := json.Marshal(map[string]interface{}{
		"description": description,
		"allowed_ips": ips,
	})
	body := bytes.Buffer{}
	body.Write(data)

	req := c.newRequest("POST", "/account/api-keys", &body)
	res, err := c.Http.Do(req)
	if err != nil {
		return nil, err
	}

	buf, err := validate(res)
	if err != nil {
		return nil, err
	}

	var model *struct {
		Attributes ApiKey `json:"attributes"`
	}
	if err = json.Unmarshal(buf, &model); err != nil {
		return nil, err
	}

	return &model.Attributes, nil
}

func (c *Client) DeleteKey(identifier string) error {
	req := c.newRequest("DELETE", "/account/api-keys/"+identifier, nil)
	res, err := c.Http.Do(req)
	if err != nil {
		return err
	}

	if _, err = validate(res); err != nil {
		return err
	}

	return nil
}
