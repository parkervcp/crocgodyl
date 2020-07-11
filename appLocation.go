package crocgodyl

import (
	"encoding/json"
	"fmt"
	"time"
)

// Application Location API

// AppLocations is the struct for all the nodes added to the panel.
// GET this from the '/api/application/locations` endpoint
type AppLocations struct {
	Object    string     `json:"object,omitempty"`
	Locations []Location `json:"data,omitempty"`
	Meta      Meta       `json:"meta,omitempty"`
}

// Location is the struct for a single location.
// GET this from the '/api/application/locations/<location_ID>` endpoint
type Location struct {
	Object     string             `json:"object,omitempty"`
	Attributes LocationAttributes `json:"attributes,omitempty"`
}

// LocationAttributes is the struct for a locations attributes.
// GET this from the '/api/application/locations/<location_ID>` endpoint
// You can only edit the short and long(description) names on a location.
type LocationAttributes struct {
	ID        int       `json:"id,omitempty"`
	Short     string    `json:"short,omitempty"`
	Long      string    `json:"long,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

// GetLocationByPage returns all available locations by page.
func (config *AppConfig) getLocationsByPage(pageID int) (locations AppLocations, err error) {
	endpoint := fmt.Sprintf("locations?page=%d", pageID)

	// Get location info from the panel
	locBytes, err := config.queryApplicationAPI(endpoint, "get", nil)
	if err != nil {
		return
	}

	// Unmarshal the bytes to a usable struct.
	err = json.Unmarshal(locBytes, &locations)
	if err != nil {
		return
	}

	return
}

// GetLocations returns all available nodes.
// Depending on how man locations you have this may take a while.
func (config *AppConfig) GetLocations() (locations AppLocations, err error) {
	// Get location info from the panel
	locBytes, err := config.queryApplicationAPI("locations", "get", nil)
	if err != nil {
		return
	}

	// Unmarshal the bytes to a usable struct.
	err = json.Unmarshal(locBytes, &locations)
	if err != nil {
		return
	}

	if locations.Meta.Pagination.TotalPages > 1 {
		for i := 1; i >= locations.Meta.Pagination.TotalPages; i++ {
			pageLocations, err := config.getLocationsByPage(i)
			if err != nil {
				return locations, err
			}
			for _, location := range pageLocations.Locations {
				locations.Locations = append(locations.Locations, location)
			}
		}
	}

	return
}

// GetLocation returns a single location by locationID.
func (config *AppConfig) GetLocation(locationID int) (location Location, err error) {
	endpoint := fmt.Sprintf("locations/%d", locationID)

	locBytes, err := config.queryApplicationAPI(endpoint, "get", nil)
	if err != nil {
		return location, err
	}

	// Get node info from the panel
	// Unmarshal the bytes to a usable struct.
	err = json.Unmarshal(locBytes, &location)
	if err != nil {
		return location, err
	}

	return location, nil
}

// CreateLocation creates a user.
func (config *AppConfig) CreateLocation(newLocation LocationAttributes) (location Location, err error) {
	newLocBytes, err := json.Marshal(newLocation)
	if err != nil {
		return location, err
	}

	// get json bytes from the panel.
	locBytes, err := config.queryApplicationAPI("locations/", "post", newLocBytes)
	if err != nil {
		return location, err
	}

	// Get server info from the panel
	// Unmarshal the bytes to a usable struct.
	err = json.Unmarshal(locBytes, &location)
	if err != nil {
		return
	}

	return
}

// EditLocation creates a user.
func (config *AppConfig) EditLocation(editLocation LocationAttributes, locationID int) (location Location, err error) {
	endpoint := fmt.Sprintf("locations/%d", locationID)

	editLocBytes, err := json.Marshal(editLocation)
	if err != nil {
		return location, err
	}

	// get json bytes from the panel.
	locBytes, err := config.queryApplicationAPI(endpoint, "patch", editLocBytes)
	if err != nil {
		return location, err
	}

	// Get server info from the panel
	// Unmarshal the bytes to a usable struct.
	err = json.Unmarshal(locBytes, &location)
	if err != nil {
		return
	}

	return
}

// DeleteLocation deletes a location.
// It only requires a locationID as an int
func (config *AppConfig) DeleteLocation(locationID int) (err error) {
	endpoint := fmt.Sprintf("locations/%d", locationID)

	// get json bytes from the panel.
	_, err = config.queryApplicationAPI(endpoint, "delete", nil)
	if err != nil {
		return err
	}

	return nil
}
