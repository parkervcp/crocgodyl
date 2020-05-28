package crocgodyl

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// Application Servers API

// AppServers is the struct for the servers on the panel.
type AppServers struct {
	Object  string   `json:"object,omitempty"`
	Servers []Server `json:"data,omitempty"`
	Meta    Meta     `json:"meta,omitempty"`
}

// Servers is the struct for a server on the panel.
type Server struct {
	Object     string           `json:"object,omitempty"`
	Attributes ServerAttributes `json:"attributes,omitempty"`
}

// ServerAttributes are the attributes for a server.
type ServerAttributes struct {
	ID            int                 `json:"id,omitempty"`
	ExternalID    interface{}         `json:"external_id,omitempty"`
	UUID          string              `json:"uuid,omitempty"`
	Identifier    string              `json:"identifier,omitempty"`
	Name          string              `json:"name,omitempty"`
	Description   string              `json:"description,omitempty"`
	Suspended     bool                `json:"suspended,omitempty"`
	Limits        ServerLimits        `json:"limits,omitempty"`
	FeatureLimits ServerFeatureLimits `json:"feature_limits,omitempty"`
	User          int                 `json:"user,omitempty"`
	Node          int                 `json:"node,omitempty"`
	Allocation    int                 `json:"allocation,omitempty"`
	Relationships ServerRelations     `json:"relationships,omitempty"`
	Nest          int                 `json:"nest,omitempty"`
	Egg           int                 `json:"egg,omitempty"`
	Pack          interface{}         `json:"pack,omitempty"`
	Container     ServerContainer     `json:"container,omitempty"`
	UpdatedAt     time.Time           `json:"updated_at,omitempty"`
	CreatedAt     time.Time           `json:"created_at,omitempty"`
}

// ServerChange is the struct for the required data for creating/modifying a server.
type ServerChange struct {
	Name          string              `json:"name,omitempty"`
	User          int                 `json:"user,omitempty"`
	Egg           int                 `json:"egg,omitempty"`
	DockerImage   string              `json:"docker_image,omitempty"`
	Startup       string              `json:"startup,omitempty"`
	Environment   map[string]string   `json:"environment,omitempty"`
	Limits        ServerLimits        `json:"limits,omitempty"`
	FeatureLimits ServerFeatureLimits `json:"feature_limits,omitempty"`
	Allocation    ServerAllocation    `json:"allocation,omitempty"`
}

// ServerLimits are the system resource limits for a server
type ServerLimits struct {
	Memory int `json:"memory,omitempty"`
	Swap   int `json:"swap,omitempty"`
	Disk   int `json:"disk,omitempty"`
	Io     int `json:"io,omitempty"`
	CPU    int `json:"cpu,omitempty"`
}

// ServerFeatureLimits this is the limit on Databases and extra Allocations on a server
type ServerFeatureLimits struct {
	Databases   int `json:"databases,omitempty"`
	Allocations int `json:"allocations,omitempty"`
}

// ServerContainer is the config on the docker container the server runs in.
type ServerContainer struct {
	StartupCommand string            `json:"startup_command,omitempty"`
	Image          string            `json:"image,omitempty"`
	Installed      bool              `json:"installed,omitempty"`
	Environment    map[string]string `json:"environment,omitempty"`
}

// ServerRelData is the data for the server relationship
type ServerRelData struct {
	Object     string                  `json:"object,omitempty"`
	Attributes ServerRelDataAttributes `json:"attributes,omitempty"`
}

// ServerRelDataAttributes are the attributes for the server relationship data
type ServerRelDataAttributes struct {
	ID       int         `json:"id,omitempty"`
	IP       string      `json:"ip,omitempty"`
	Alias    interface{} `json:"alias,omitempty"`
	Port     int         `json:"port,omitempty"`
	Assigned bool        `json:"assigned,omitempty"`
}

// ServerAllocation is only used when creating a server
type ServerAllocation struct {
	Default    int   `json:"default,omitempty"`
	Additional []int `json:"additional,omitempty"`
}

// ServerRelations is the struct for Relationships for a Servers
type ServerRelations struct {
	Allocations struct {
		Object string                    `json:"object,omitempty"`
		Data   []ServerAllocRelationData `json:"data,omitempty"`
	} `json:"allocations,omitempty"`
}

// ServerAllocRelationData is the struct for Allocation Relationships on a Servers
type ServerAllocRelationData struct {
	Object     string                 `json:"object,omitempty"`
	Attributes []SererAllocAttributes `json:"data,omitempty"`
}

type SererAllocAttributes struct {
	ID       int    `json:"id,omitempty"`
	IP       string `json:"ip,omitempty"`
	Alias    string `json:"alias,omitempty"`
	Port     int    `json:"port,omitempty"`
	Assigned bool   `json:"assigned,omitempty"`
}

// getServersByPage returns all available locations by page.
func (config *CrocConfig) getServersByPage(pageID int) (servers AppServers, err error) {
	// Get location info from the panel
	serverBytes, err := config.queryApplicationAPI(fmt.Sprintf("servers?page=%d", pageID), "get", nil)
	if err != nil {
		return
	}

	// Unmarshal the bytes to a usable struct.
	err = json.Unmarshal(serverBytes, &servers)
	if err != nil {
		return
	}

	return
}

// GetServers returns all available servers.
func (config *CrocConfig) GetServers() (servers AppServers, err error) {
	// Get server info from the panel
	serverBytes, err := config.queryApplicationAPI("servers", "get", nil)
	if err != nil {
		return
	}

	// Unmarshal the bytes to a usable struct.
	err = json.Unmarshal(serverBytes, &servers)
	if err != nil {
		return
	}

	for i := 1; i >= servers.Meta.Pagination.TotalPages; i++ {
		pageServers, err := config.getServersByPage(i)
		if err != nil {
			return servers, err
		}
		for _, server := range pageServers.Servers {
			servers.Servers = append(servers.Servers, server)
		}
	}

	return
}

// GetServer returns Information on a single server.
func (config *CrocConfig) GetServer(serverID int) (server Server, err error) {
	// get json bytes from the panel.
	fmt.Sprintf("servers/%d?include=allocations", serverID)
	serverBytes, err := config.queryApplicationAPI(fmt.Sprintf("servers/%d?include=allocations", serverID), "get", nil)
	if err != nil {
		return
	}

	// Get server info from the panel
	// Unmarshal the bytes to a usable struct.
	err = json.Unmarshal(serverBytes, &server)
	if err != nil {
		return
	}

	return
}

// GetServerAllocations will return an array of SererAllocAttributes
func (config *CrocConfig) GetServerAllocations(serverID int) (serverAllocations []SererAllocAttributes, err error) {
	var server Server

	// get json bytes from the panel.
	serverAllocBytes, err := config.queryApplicationAPI(fmt.Sprintf("servers/%d?include=allocations", serverID), "get", nil)
	if err != nil {
		return serverAllocations, err
	}

	// Get server info from the panel
	// Unmarshal the bytes to a usable struct.
	err = json.Unmarshal(serverAllocBytes, &server)
	if err != nil {
		return
	}

	for i, port := range server.Attributes.Relationships.Allocations.Data {
		serverAllocations = append(serverAllocations, port.Attributes[i])
	}

	return
}

// CreateServer creates a new server via the API.
// A complete ServerChange is required.
func (config *CrocConfig) CreateServer(newServer ServerChange) (server Server, err error) {
	newServerBytes, err := json.Marshal(newServer)
	if err != nil {
		return
	}

	// get json bytes from the panel.
	serverBytes, err := config.queryApplicationAPI("servers", "post", newServerBytes)
	if err != nil {
		return
	}

	// Get server info from the panel
	// Unmarshal the bytes to a usable struct.
	err = json.Unmarshal(serverBytes, &server)
	if err != nil {
		return
	}

	return
}

// EditServerDetails creates a new server via the API.
// The server name and user are required when updating a server.
func (config *CrocConfig) EditServerDetails(newServer ServerChange, serverID int) (server Server, err error) {
	editServerBytes, err := json.Marshal(newServer)
	if err != nil {
		return
	}

	// get json bytes from the panel.
	serverBytes, err := config.queryApplicationAPI("servers/"+strconv.Itoa(serverID)+"/details", "patch", editServerBytes)
	if err != nil {
		return
	}

	// Get server info from the panel
	// Unmarshal the bytes to a usable struct.
	err = json.Unmarshal(serverBytes, &server)
	if err != nil {
		return
	}

	return
}

// EditServerBuild creates a new server via the API.
// The server name and user are required when updating a server.
func (config *CrocConfig) EditServerBuild(newServer ServerChange, serverID int) (server Server, err error) {
	editServerBytes, err := json.Marshal(newServer)
	if err != nil {
		return
	}

	// get json bytes from the panel.
	serverBytes, err := config.queryApplicationAPI("servers/"+strconv.Itoa(serverID)+"/build", "patch", editServerBytes)
	if err != nil {
		return
	}

	// Get server info from the panel
	// Unmarshal the bytes to a usable struct.
	err = json.Unmarshal(serverBytes, &server)
	if err != nil {
		return
	}

	return
}

// EditServerStartup creates a new server via the API.
// The server name and user are required when updating a server.
func (config *CrocConfig) EditServerStartup(newServer ServerChange, serverID int) (server Server, err error) {
	editServerBytes, err := json.Marshal(newServer)
	if err != nil {
		return
	}

	// get json bytes from the panel.
	serverBytes, err := config.queryApplicationAPI("servers/"+strconv.Itoa(serverID)+"/startup", "patch", editServerBytes)
	if err != nil {
		return
	}

	// Get server info from the panel
	// Unmarshal the bytes to a usable struct.
	err = json.Unmarshal(serverBytes, &server)
	if err != nil {
		return
	}

	return
}

// DeleteServer deletes a server.
// It only requires a server id as a string
func (config *CrocConfig) DeleteServer(serverID int) error {
	// get json bytes from the panel.
	_, err := config.queryApplicationAPI("servers/"+strconv.Itoa(serverID), "delete", nil)
	if err != nil {
		return err
	}

	return nil
}
