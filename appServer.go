package crocgodyl

import (
	"encoding/json"
	"time"
)

// Application Server API

// Servers is the struct for the servers on the panel.
type Servers struct {
	Object string   `json:"object"`
	Server []Server `json:"data"`
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

// Server is the struct for a server on the panel.
type Server struct {
	Object     string           `json:"object"`
	Attributes ServerAttributes `json:"attributes"`
}

// ServerAttributes are the attributes for a server.
type ServerAttributes struct {
	ID            int                 `json:"id"`
	ExternalID    interface{}         `json:"external_id"`
	UUID          string              `json:"uuid"`
	Identifier    string              `json:"identifier"`
	Name          string              `json:"name"`
	Description   string              `json:"description"`
	Suspended     bool                `json:"suspended"`
	Limits        ServerLimits        `json:"limits"`
	FeatureLimits ServerFeatureLimits `json:"feature_limits"`
	User          int                 `json:"user"`
	Node          int                 `json:"node"`
	Allocation    int                 `json:"allocation"`
	Nest          int                 `json:"nest"`
	Egg           int                 `json:"egg"`
	Pack          interface{}         `json:"pack"`
	Container     ServerContainer     `json:"container"`
	UpdatedAt     time.Time           `json:"updated_at"`
	CreatedAt     time.Time           `json:"created_at"`
}

// ServerChange is the struct for the required data for creating/modifying a server.
type ServerChange struct {
	Name          string              `json:"name"`
	User          int                 `json:"user"`
	Egg           int                 `json:"egg"`
	DockerImage   string              `json:"docker_image"`
	Startup       string              `json:"startup"`
	Environment   map[string]string   `json:"environment"`
	Limits        ServerLimits        `json:"limits"`
	FeatureLimits ServerFeatureLimits `json:"feature_limits"`
	Allocation    ServerAllocation    `json:"allocation"`
}

// ServerLimits are the system resource limits for a server
type ServerLimits struct {
	Memory int `json:"memory"`
	Swap   int `json:"swap"`
	Disk   int `json:"disk"`
	Io     int `json:"io"`
	CPU    int `json:"cpu"`
}

// ServerFeatureLimits this is the limit on Databases and extra Allocations on a server
type ServerFeatureLimits struct {
	Databases   int `json:"databases"`
	Allocations int `json:"allocations"`
}

// ServerContainer is the config on the docker container the server runs in.
type ServerContainer struct {
	StartupCommand string            `json:"startup_command"`
	Image          string            `json:"image"`
	Installed      bool              `json:"installed"`
	Environment    map[string]string `json:"environment"`
}

// ServerAllocation is only used when creating a server
type ServerAllocation struct {
	Default int `json:"default"`
}

// GetServers returns all available servers.
func GetServers() (Servers, error) {
	var servers Servers

	// get json bytes from the panel.
	sbytes, err := queryPanelAPI("servers", "get", nil)
	if err != nil {
		return servers, err
	}

	// Get server info from the panel
	// Unmarshal the bytes to a usable struct.
	err = json.Unmarshal(sbytes, &servers)
	if err != nil {
		return servers, err
	}

	return servers, nil
}

// CreateServer creates a new server via the API.
// A complete ServerChange is required. See the example at
func CreateServer(newServer ServerChange) (Server, error) {
	var serverDetails Server

	nsbytes, err := json.Marshal(newServer)
	if err != nil {
		return serverDetails, err
	}

	// get json bytes from the panel.
	sbytes, err := queryPanelAPI("servers", "post", nsbytes)
	if err != nil {
		return serverDetails, err
	}

	// Get server info from the panel
	// Unmarshal the bytes to a usable struct.
	err = json.Unmarshal(sbytes, &serverDetails)
	if err != nil {
		return serverDetails, err
	}

	return serverDetails, nil
}

// DeleteServer deletes a server.
// It only requires a server id as a string
func DeleteServer(serverid int) error {
	// get json bytes from the panel.
	_, err := queryPanelAPI("servers/"+string(serverid), "delete", nil)
	if err != nil {
		return err
	}

	return nil
}
