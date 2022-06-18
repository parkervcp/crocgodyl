package crocgodyl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

type User struct {
	ID         int64      `json:"id"`
	ExternalID string     `json:"external_id"`
	UUID       string     `json:"uuid"`
	Username   string     `json:"username"`
	Email      string     `json:"email"`
	FirstName  string     `json:"first_name"`
	LastName   string     `json:"last_name"`
	Language   string     `json:"language"`
	RootAdmin  bool       `json:"root_admin"`
	TwoFactor  bool       `json:"2fa"`
	CreatedAt  *time.Time `json:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at,omitempty"`
}

func (u *User) FullName() string {
	return u.FirstName + " " + u.LastName
}

func (a *Application) Users() ([]*User, error) {
	req := a.newRequest("GET", "/users", nil)
	res, err := a.Http.Do(req)
	if err != nil {
		return nil, err
	}

	buf, err := validate(res)
	if err != nil {
		return nil, err
	}

	var model struct {
		Data []struct {
			Attributes *User `json:"attributes"`
		} `json:"data"`
	}
	if err = json.Unmarshal(buf, &model); err != nil {
		return nil, err
	}

	users := make([]*User, 0, len(model.Data))
	for _, u := range model.Data {
		users = append(users, u.Attributes)
	}

	return users, nil
}

func (a *Application) User(id int64) (*User, error) {
	req := a.newRequest("GET", fmt.Sprintf("/users/%d", id), nil)
	res, err := a.Http.Do(req)
	if err != nil {
		return nil, err
	}

	buf, err := validate(res)
	if err != nil {
		return nil, err
	}

	var model struct {
		Attributes User `json:"attributes"`
	}
	if err = json.Unmarshal(buf, &model); err != nil {
		return nil, err
	}

	return &model.Attributes, nil
}

func (a *Application) UserExternal(id string) (*User, error) {
	req := a.newRequest("GET", fmt.Sprintf("/users/external/%s", id), nil)
	res, err := a.Http.Do(req)
	if err != nil {
		return nil, err
	}

	buf, err := validate(res)
	if err != nil {
		return nil, err
	}

	var model struct {
		Attributes User `json:"attributes"`
	}
	if err = json.Unmarshal(buf, &model); err != nil {
		return nil, err
	}

	return &model.Attributes, nil
}

type CreateUserDescriptor struct {
	ExternalID string `json:"external_id,omitempty"`
	Email      string `json:"email"`
	Username   string `json:"username"`
	Password   string `json:"password,omitempty"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Language   string `json:"language,omitempty"`
	RootAdmin  bool   `json:"root_admin,omitempty"`
}

func (a *Application) CreateUser(fields CreateUserDescriptor) (*User, error) {
	data, _ := json.Marshal(fields)
	body := bytes.Buffer{}
	body.Write(data)

	req := a.newRequest("POST", "/users", &body)
	res, err := a.Http.Do(req)
	if err != nil {
		return nil, err
	}

	buf, err := validate(res)
	if err != nil {
		return nil, err
	}

	var model struct {
		Attributes User `json:"attributes"`
	}
	if err = json.Unmarshal(buf, &model); err != nil {
		return nil, err
	}

	return &model.Attributes, nil
}

type UpdateUserDescriptor struct {
	ExternalID string `json:"external_id,omitempty"`
	Email      string `json:"email,omitempty"`
	Username   string `json:"username,omitempty"`
	Password   string `json:"password,omitempty"`
	FirstName  string `json:"first_name,omitempty"`
	LastName   string `json:"last_name,omitempty"`
	Language   string `json:"language,omitempty"`
	RootAdmin  bool   `json:"root_admin,omitempty"`
}

func (a *Application) UpdateUser(id int64, fields UpdateUserDescriptor) (*User, error) {
	data, _ := json.Marshal(fields)
	body := bytes.Buffer{}
	body.Write(data)

	req := a.newRequest("PATCH", fmt.Sprintf("/users/%d", id), &body)
	res, err := a.Http.Do(req)
	if err != nil {
		return nil, err
	}

	buf, err := validate(res)
	if err != nil {
		return nil, err
	}

	var model struct {
		Attributes User `json:"attributes"`
	}
	if err = json.Unmarshal(buf, &model); err != nil {
		return nil, err
	}

	return &model.Attributes, nil
}

func (a *Application) DeleteUser(id int64) error {
	req := a.newRequest("DELETE", fmt.Sprintf("/users/%d", id), nil)
	res, err := a.Http.Do(req)
	if err != nil {
		return err
	}

	_, err = validate(res)
	return err
}
