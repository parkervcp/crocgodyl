package crocgodyl

import (
	"encoding/json"
	"fmt"
	"time"
)

// Application Users API

// AppUsers is the struct for all the panel users.
// GET this from the '/api/application/users` endpoint
type AppUsers struct {
	Object string `json:"object,omitempty"`
	Users  []User `json:"data,omitempty"`
	Meta   Meta   `json:"meta,omitempty"`
}

// Users is the struct for all the panel users.
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

// GetUserByPage returns Information on users by their page number.
// The externalID is a string as that is what the panel requires.
func (config *CrocConfig) getUserByPage(pageID int) (users AppUsers, err error) {
	// get json bytes from the panel.
	userBytes, err := config.queryApplicationAPI(fmt.Sprintf("users/%d", pageID), "get", nil)
	if err != nil {
		return users, err
	}

	// Unmarshal the bytes to a usable struct.
	err = json.Unmarshal(userBytes, &users)
	if err != nil {
		return users, err
	}

	return users, nil
}

// GetUsers returns Information on all users.
func (config *CrocConfig) GetUsers() (users AppUsers, err error) {
	// Get location info from the panel
	userBytes, err := config.queryApplicationAPI("users", "get", nil)
	if err != nil {
		return
	}

	// Unmarshal the bytes to a usable struct.
	err = json.Unmarshal(userBytes, &users)
	if err != nil {
		return
	}

	for i := 1; i >= users.Meta.Pagination.TotalPages; i++ {
		pageUsers, err := config.getUserByPage(i)
		if err != nil {
			return users, err
		}
		for _, user := range pageUsers.Users {
			users.Users = append(users.Users, user)
		}
	}

	return
}

// GetUser returns Information on a single user.
func (config *CrocConfig) GetUser(userID int) (user User, err error) {
	// get json bytes from the panel.
	userBytes, err := config.queryApplicationAPI(fmt.Sprintf("users/%d", userID), "get", nil)
	if err != nil {
		return
	}

	// Unmarshal the bytes to a usable struct.
	err = json.Unmarshal(userBytes, &user)
	if err != nil {
		return
	}

	return
}

// GetUserByExternal returns Information on a single user by their externalID.
// The externalID is a string as that is what the panel requires.
func (config *CrocConfig) GetUserByExternal(externalID string) (user User, err error) {
	// get json bytes from the panel.
	userBytes, err := config.queryApplicationAPI(fmt.Sprintf("users/%s", externalID), "get", nil)
	if err != nil {
		return
	}

	// Unmarshal the bytes to a usable struct.
	err = json.Unmarshal(userBytes, &user)
	if err != nil {
		return
	}

	return
}

// CreateUser creates a user.
func (config *CrocConfig) CreateUser(newUser UserAttributes) (user User, err error) {
	newUserBytes, err := json.Marshal(newUser)
	if err != nil {
		return
	}

	// get json bytes from the panel.
	userBytes, err := config.queryApplicationAPI("users", "post", newUserBytes)
	if err != nil {
		return
	}

	// Get user info from the panel
	// Unmarshal the bytes to a usable struct.
	err = json.Unmarshal(userBytes, &user)
	if err != nil {
		return
	}

	return
}

// EditUser creates a user.
// Send a UserAttributes to the panel to update the user.
// You cannot edit the id or created/updated fields for the user.
func (config *CrocConfig) EditUser(editUser UserAttributes, userID int) (user User, err error) {
	editUserBytes, err := json.Marshal(editUser)
	if err != nil {
		return
	}

	// get json bytes from the panel.
	userBytes, err := config.queryApplicationAPI(fmt.Sprintf("users/%d", userID), "patch", editUserBytes)
	if err != nil {
		return
	}

	// Get user info from the panel
	// Unmarshal the bytes to a usable struct.
	err = json.Unmarshal(userBytes, &user)
	if err != nil {
		return
	}

	return
}

// DeleteUser deletes a user.
// It only requires a user id as a string
func (config *CrocConfig) DeleteUser(userID int) (err error) {
	// get json bytes from the panel.
	_, err = config.queryApplicationAPI(fmt.Sprintf("users/%d", userID), "delete", nil)
	if err != nil {
		return
	}

	return
}
