package crocgodyl

import (
	"encoding/json"
	"fmt"
	"time"
)

// Application User API

// Users is the struct for all the panel users.
// GET this from the '/api/application/users` endpoint
type Users struct {
	Object string `json:"object,omitempty"`
	User   []User `json:"data,omitempty"`
	Meta   Meta   `json:"meta,omitempty"`
}

// User is the struct for all the panel users.
// GET this from the '/api/application/users/<user_ID>` endpoint
// The panel does not and will not return a password.
// You can update a password using the API.
type User struct {
	Object     string         `json:"object,omitempty"`
	Attributes UserAttributes `json:"attributes,omitempty"`
}

// UserAttributes is the struct for all the panel users.
type UserAttributes struct {
	ID         int       `json:"id,omitempty"`
	ExternalID string    `json:"external_id,omitempty"`
	UUID       string    `json:"uuid,omitempty"`
	Username   string    `json:"username,omitempty"`
	Email      string    `json:"email,omitempty"`
	FirstName  string    `json:"first_name,omitempty"`
	LastName   string    `json:"last_name,omitempty"`
	Language   string    `json:"language,omitempty"`
	RootAdmin  bool      `json:"root_admin,omitempty"`
	TwoFa      bool      `json:"2fa,omitempty"`
	Password   string    `json:"password,omitempty"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
	UpdatedAt  time.Time `json:"updated_at,omitempty"`
}

// GetUsers returns information on all users.
func (config *CrocConfig) GetUsers() (Users, error) {
	var users Users

	// get json bytes from the panel.
	uBytes, err := config.queryPanelAPI("users", "get", nil)
	if err != nil {
		return users, err
	}

	// Unmarshal the bytes to a usable struct.
	err = json.Unmarshal(uBytes, &users)
	if err != nil {
		return users, err
	}

	return users, nil
}

// GetUser returns information on a single user.
func (config *CrocConfig) GetUser(userID int) (User, error) {
	var user User
	endpoint := fmt.Sprintf("users/%d", userID)

	// get json bytes from the panel.
	uBytes, err := config.queryPanelAPI(endpoint, "get", nil)
	if err != nil {
		return user, err
	}

	// Unmarshal the bytes to a usable struct.
	err = json.Unmarshal(uBytes, &user)
	if err != nil {
		return user, err
	}

	return user, nil
}

// GetUserByExternal returns information on a single user by their external id.
func (config *CrocConfig) GetUserByExternal(externalID string) (User, error) {
	var user User
	endpoint := fmt.Sprintf("users/%s", externalID)

	// get json bytes from the panel.
	uBytes, err := config.queryPanelAPI(endpoint, "get", nil)
	if err != nil {
		return user, err
	}

	// Unmarshal the bytes to a usable struct.
	err = json.Unmarshal(uBytes, &user)
	if err != nil {
		return user, err
	}

	return user, nil
}

// GetUserByPage returns information on users by their page number.
func (config *CrocConfig) GetUserByPage(pageID int) (User, error) {
	var user User
	endpoint := fmt.Sprintf("users/%d", pageID)

	// get json bytes from the panel.
	uBytes, err := config.queryPanelAPI(endpoint, "get", nil)
	if err != nil {
		return user, err
	}

	// Unmarshal the bytes to a usable struct.
	err = json.Unmarshal(uBytes, &user)
	if err != nil {
		return user, err
	}

	return user, nil
}

// CreateUser creates a user.
func (config *CrocConfig) CreateUser(newUser UserAttributes) (User, error) {
	var userDetails User

	nUBytes, err := json.Marshal(newUser)
	if err != nil {
		return userDetails, err
	}

	// get json bytes from the panel.
	uBytes, err := config.queryPanelAPI("users", "post", nUBytes)
	if err != nil {
		return userDetails, err
	}

	// Get user info from the panel
	// Unmarshal the bytes to a usable struct.
	err = json.Unmarshal(uBytes, &userDetails)
	if err != nil {
		return userDetails, err
	}

	return userDetails, nil
}

// EditUser edits the information of a specified user.
func (config *CrocConfig) EditUser(editUser UserAttributes, userID int) (User, error) {
	var userDetails User
	endpoint := fmt.Sprintf("users/%d", userID)

	eUBytes, err := json.Marshal(editUser)
	if err != nil {
		return userDetails, err
	}

	// get json bytes from the panel.
	uBytes, err := config.queryPanelAPI(endpoint, "patch", eUBytes)
	if err != nil {
		return userDetails, err
	}

	// Get user info from the panel
	// Unmarshal the bytes to a usable struct.
	err = json.Unmarshal(uBytes, &userDetails)
	if err != nil {
		return userDetails, err
	}

	return userDetails, nil
}

// DeleteUser deletes a user.
func (config *CrocConfig) DeleteUser(userID int) error {
	endpoint := fmt.Sprintf("users/%d", userID)

	// get json bytes from the panel.
	_, err := config.queryPanelAPI(endpoint, "delete", nil)
	if err != nil {
		return err
	}

	return nil
}
