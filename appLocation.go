package crocgodyl

import (
	"encoding/json"
	"fmt"
	"time"
)

// Application Location API

// Locations is the struct for all the nodes added to the panel.
// GET this from the '/api/application/locations` endpoint
type Locations struct {
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

// GetLocations returns all available nodes.
// Depending on how man locations you have this may take a while.
func (config *CrocConfig) GetLocations() (Locations, error) {
	var locations Locations
	var locationsAll Locations

	pages, err := config.GetLocationByPage(1)
	if err != nil {
		return locations, err
	}

	for i := 1; i >= pages.Meta.Pagination.TotalPages; i++ {
		locations, err := config.GetLocationByPage(i)
		if err != nil {
			return locations, err
		}
		for _, location := range locations.Locations {
			locationsAll.Locations = append(locationsAll.Locations, location)
		}
	}

	return locationsAll, nil
}

// GetLocation returns a single location by locationID.
func (config *CrocConfig) GetLocation(locationID int) (Location, error) {
	var location Location
	endpoint := fmt.Sprintf("locations/%d", locationID)

	locBytes, err := config.queryPanelAPI(endpoint, "get", nil)
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

// GetLocationByPage returns all available locations by page.
func (config *CrocConfig) GetLocationByPage(pageID int) (Locations, error) {
	var locations Locations
	endpoint := fmt.Sprintf("locations?page=%d", pageID)

	locBytes, err := config.queryPanelAPI(endpoint, "get", nil)
	if err != nil {
		return locations, err
	}

	// Get node info from the panel
	// Unmarshal the bytes to a usable struct.
	err = json.Unmarshal(locBytes, &locations)
	if err != nil {
		return locations, err
	}

	return locations, nil
}

// CreateLocation creates a user.
func (config *CrocConfig) CreateLocation(newLocation LocationAttributes) (Location, error) {
	var locationDetails Location

	newLocBytes, err := json.Marshal(newLocation)
	if err != nil {
		return locationDetails, err
	}

	// get json bytes from the panel.
	locBytes, err := config.queryPanelAPI("locations/", "post", newLocBytes)
	if err != nil {
		return locationDetails, err
	}

	// Get server info from the panel
	// Unmarshal the bytes to a usable struct.
	err = json.Unmarshal(locBytes, &locationDetails)
	if err != nil {
		return locationDetails, err
	}

	return locationDetails, nil
}

// EditLocation creates a user.
func (config *CrocConfig) EditLocation(editLocation LocationAttributes, locationID int) (Location, error) {
	var locationDetails Location
	endpoint := fmt.Sprintf("locations/%d", locationID)

	editLocBytes, err := json.Marshal(editLocation)
	if err != nil {
		return locationDetails, err
	}

	// get json bytes from the panel.
	locBytes, err := config.queryPanelAPI(endpoint, "patch", editLocBytes)
	if err != nil {
		return locationDetails, err
	}

	// Get server info from the panel
	// Unmarshal the bytes to a usable struct.
	err = json.Unmarshal(locBytes, &locationDetails)
	if err != nil {
		return locationDetails, err
	}

	return locationDetails, nil
}

// DeleteLocation deletes a location.
// It only requires a locationID as an int
func (config *CrocConfig) DeleteLocation(locationID int) error {
	endpoint := fmt.Sprintf("locations/%d", locationID)

	// get json bytes from the panel.
	_, err := config.queryPanelAPI(endpoint, "delete", nil)
	if err != nil {
		return err
	}

	return nil
}
