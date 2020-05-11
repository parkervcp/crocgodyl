package crocgodyl

import (
	"encoding/json"
	"strconv"
	"time"
)

// Application Server API

// Servers is the struct for the servers on the panel.
type Servers struct {
	Object string   `json:"object,omitempty"`
	Server []Server `json:"data,omitempty"`
	Meta   Meta     `json:"meta,omitempty"`
}

// Server is the struct for a server on the panel.
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

// ServerRelations is the struct for Relationships for a Server
type ServerRelations struct {
	Allocations struct {
		Object string                    `json:"object,omitempty"`
		Data   []ServerAllocRelationData `json:"data,omitempty"`
	} `json:"allocations,omitempty"`
}

// ServerAllocRelationData is the struct for Allocation Relationships on a Server
type ServerAllocRelationData struct {
	Object     string `json:"object,omitempty"`
	Attributes struct {
		ID       int    `json:"id,omitempty"`
		IP       string `json:"ip,omitempty"`
		Alias    string `json:"alias,omitempty"`
		Port     int    `json:"port,omitempty"`
		Assigned bool   `json:"assigned,omitempty"`
	}
}

// GetServers returns all available servers.
func (config *CrocConfig) GetServers() (Servers, error) {
	var servers Servers

	// get json bytes from the panel.
	sBytes, err := config.queryPanelAPI("servers", "get", nil)
	if err != nil {
		return servers, err
	}

	// Get server info from the panel
	// Unmarshal the bytes to a usable struct.
	err = json.Unmarshal(sBytes, &servers)
	if err != nil {
		return servers, err
	}

	return servers, nil
}

// GetServer returns Information on a single server.
func (config *CrocConfig) GetServer(serverId int) (Server, error) {
	var server Server

	// get json bytes from the panel.
	sBytes, err := config.queryPanelAPI("servers/"+strconv.Itoa(serverId)+"?include=allocations", "get", nil)
	if err != nil {
		return server, err
	}

	// Get server info from the panel
	// Unmarshal the bytes to a usable struct.
	err = json.Unmarshal(sBytes, &server)
	if err != nil {
		return server, err
	}

	return server, nil
}

// GetServerAllocations will return a list of allocations for the server in a []int array
func (config *CrocConfig) GetServerAllocations(serverId int) ([]int, error) {
	var allServerAlloc []int

	// get json bytes from the panel.
	svAllocBytes, err := config.queryPanelAPI("servers/"+strconv.Itoa(serverId)+"?include=allocations", "get", nil)
	if err != nil {
		return allServerAlloc, err
	}

	// Get server info from the panel
	// Unmarshal the bytes to a usable struct.
	err = json.Unmarshal(svAllocBytes, &allServerAlloc)
	if err != nil {
		return allServerAlloc, err
	}

	return allServerAlloc, nil
}

// CreateServer creates a new server via the API.
// A complete ServerChange is required.
func (config *CrocConfig) CreateServer(newServer ServerChange) (Server, error) {
	var serverDetails Server

	nsBytes, err := json.Marshal(newServer)
	if err != nil {
		return serverDetails, err
	}

	// get json bytes from the panel.
	sBytes, err := config.queryPanelAPI("servers", "post", nsBytes)
	if err != nil {
		return serverDetails, err
	}

	// Get server info from the panel
	// Unmarshal the bytes to a usable struct.
	err = json.Unmarshal(sBytes, &serverDetails)
	if err != nil {
		return serverDetails, err
	}

	return serverDetails, nil
}

// EditServerDetails edits the server details for the specified server.
// The server name and user are required when updating a server.
func (config *CrocConfig) EditServerDetails(newServer ServerChange, serverid int) (Server, error) {
	var serverDetails Server

	esBytes, err := json.Marshal(newServer)
	if err != nil {
		return serverDetails, err
	}

	// get json bytes from the panel.
	sBytes, err := config.queryPanelAPI("servers/"+strconv.Itoa(serverid)+"/details", "patch", esBytes)
	if err != nil {
		return serverDetails, err
	}

	// Get server info from the panel
	// Unmarshal the bytes to a usable struct.
	err = json.Unmarshal(sBytes, &serverDetails)
	if err != nil {
		return serverDetails, err
	}

	return serverDetails, nil
}

// EditServerBuild edits the build information for the specified server.
// The server name and user are required when updating a server.
func (config *CrocConfig) EditServerBuild(newServer ServerChange, serverid int) (Server, error) {
	var serverDetails Server

	esBytes, err := json.Marshal(newServer)
	if err != nil {
		return serverDetails, err
	}

	// get json bytes from the panel.
	sBytes, err := config.queryPanelAPI("servers/"+strconv.Itoa(serverid)+"/build", "patch", esBytes)
	if err != nil {
		return serverDetails, err
	}

	// Get server info from the panel
	// Unmarshal the bytes to a usable struct.
	err = json.Unmarshal(sBytes, &serverDetails)
	if err != nil {
		return serverDetails, err
	}

	return serverDetails, nil
}

// EditServerStartup edits the startup information for the specified server.
// The server name and user are required when updating a server.
func (config *CrocConfig) EditServerStartup(newServer ServerChange, serverid int) (Server, error) {
	var serverDetails Server

	esBytes, err := json.Marshal(newServer)
	if err != nil {
		return serverDetails, err
	}

	// get json bytes from the panel.
	sBytes, err := config.queryPanelAPI("servers/"+strconv.Itoa(serverid)+"/startup", "patch", esBytes)
	if err != nil {
		return serverDetails, err
	}

	// Get server info from the panel
	// Unmarshal the bytes to a usable struct.
	err = json.Unmarshal(sBytes, &serverDetails)
	if err != nil {
		return serverDetails, err
	}

	return serverDetails, nil
}

// DeleteServer deletes a server.
// It only requires a server id as a string
func (config *CrocConfig) DeleteServer(serverId int) error {
	// get json bytes from the panel.
	_, err := config.queryPanelAPI("servers/"+strconv.Itoa(serverId), "delete", nil)
	if err != nil {
		return err
	}

	return nil
}
