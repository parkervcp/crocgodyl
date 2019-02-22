package crocgodyl

import (
	"encoding/json"
	"time"
)

// Application Location API

// Locations is the struct for all the nodes added to the panel.
// GET this from the '/api/application/locations` endpoint
type Locations struct {
	Object   string     `json:"object"`
	Location []Location `json:"data"`
	Meta     struct {
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

// Location is the struct for a single location.
// GET this from the '/api/application/locations/<location_ID>` endpoint
type Location struct {
	Object     string `json:"object"`
	Attributes struct {
		ID        int       `json:"id"`
		Short     string    `json:"short"`
		Long      string    `json:"long"`
		UpdatedAt time.Time `json:"updated_at"`
		CreatedAt time.Time `json:"created_at"`
	} `json:"attributes"`
}

// LocationEdit is the struct for the json when editing/creating a location.
// GET this from the '/api/application/locations/<location_ID>` endpoint
type LocationEdit struct {
	Short string `json:"short"`
	Long  string `json:"long"`
}

// LocatioCreateResponse is the struct for the response when creating a location.
// GET this from the '/api/application/locations/<location_ID>` endpoint
type LocatioCreateResponse struct {
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

// GetLocations returns all available nodes.
func GetLocations() (Locations, error) {
	var locations Locations

	lbytes, err := queryPanelAPI("locations", "get", nil)
	if err != nil {
		return locations, err
	}

	// Get node info from the panel
	// Unmarshal the bytes to a usable struct.
	err = json.Unmarshal(lbytes, &locations)
	if err != nil {
		return locations, err
	}

	return locations, nil
}
