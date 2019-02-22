package crocgodyl

import (
	"encoding/json"
	"time"
)

// Application User API

// Users is the struct for all the panel users.
// GET this from the '/api/application/users` endpoint
type Users struct {
	Object string `json:"object"`
	User   []User `json:"data"`
	Meta   struct {
		Pagination struct {
			Total       int           `json:"total"`
			Count       int           `json:"count"`
			PerPage     int           `json:"per_page"`
			CurrentPage int           `json:"current_page"`
			TotalPages  int           `json:"total_pages"`
			Links       []interface{} `json:"links"`
		} `json:"pagination"`
	} `json:"meta"`
}

// User is the struct for all the panel users.
// GET this from the '/api/application/users/<user_ID>` endpoint
type User struct {
	Object     string `json:"object"`
	Attributes struct {
		ID         int         `json:"id"`
		ExternalID interface{} `json:"external_id"`
		UUID       string      `json:"uuid"`
		Username   string      `json:"username"`
		Email      string      `json:"email"`
		FirstName  string      `json:"first_name"`
		LastName   string      `json:"last_name"`
		Language   string      `json:"language"`
		RootAdmin  bool        `json:"root_admin"`
		TwoFa      bool        `json:"2fa"`
		CreatedAt  time.Time   `json:"created_at"`
		UpdatedAt  time.Time   `json:"updated_at"`
	} `json:"attributes"`
}

// PanelUserEdit is the struct for creating a panel user.
// POST this to the '/api/application/users/` endpoint
// PATCH this to the '/api/application/users/<userID>` endpoint
type PanelUserEdit struct {
	Username   string `json:"username"`
	Email      string `json:"email"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Language   string `json:"language"`
	ExternalID string `json:"external_id"`
	RootAdmin  bool   `json:"root_admin"`
}

// UserCreateResponse is the struct for the response when creating a user.
// POST this from the '/api/application/users/<user_ID>` endpoint
type UserCreateResponse struct {
	Object     string `json:"object"`
	Attributes struct {
		ID         int         `json:"id"`
		ExternalID interface{} `json:"external_id"`
		UUID       string      `json:"uuid"`
		Username   string      `json:"username"`
		Email      string      `json:"email"`
		FirstName  string      `json:"first_name"`
		LastName   string      `json:"last_name"`
		Language   string      `json:"language"`
		RootAdmin  bool        `json:"root_admin"`
		TwoFa      bool        `json:"2fa"`
		CreatedAt  time.Time   `json:"created_at"`
		UpdatedAt  time.Time   `json:"updated_at"`
	} `json:"attributes"`
	Meta struct {
		Resource string `json:"resource"`
	} `json:"meta"`
}

// UserUpdateResponse is the struct for the response when editing a user.
// PATCH this from the '/api/application/users/<user_ID>` endpoint
type UserUpdateResponse struct {
	Object     string `json:"object"`
	Attributes struct {
		ID         int       `json:"id"`
		ExternalID string    `json:"external_id"`
		UUID       string    `json:"uuid"`
		Username   string    `json:"username"`
		Email      string    `json:"email"`
		FirstName  string    `json:"first_name"`
		LastName   string    `json:"last_name"`
		Language   string    `json:"language"`
		RootAdmin  bool      `json:"root_admin"`
		TwoFa      bool      `json:"2fa"`
		CreatedAt  time.Time `json:"created_at"`
		UpdatedAt  time.Time `json:"updated_at"`
	} `json:"attributes"`
}

// GetUsers returns all available nodes.
func GetUsers() (Users, error) {
	var users Users

	// get json bytes from the panel.
	ubytes, err := queryPanelAPI("users", "get", nil)
	if err != nil {
		return users, err
	}

	// Unmarshal the bytes to a usable struct.
	err = json.Unmarshal(ubytes, &users)
	if err != nil {
		return users, err
	}

	return users, nil
}

// SetUser changes user settings
func SetUser() {

}
