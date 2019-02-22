package crocgodyl

import (
	"encoding/json"
	"time"
)

// Application Node API

// Nodes is the struct for all the nodes added to the panel.
type Nodes struct {
	Object string `json:"object"`
	Node   []Node `json:"data"`
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

// Node is the struct for a single node.
type Node struct {
	Object     string `json:"object"`
	Attributes struct {
		ID                 int       `json:"id"`
		Public             bool      `json:"public"`
		Name               string    `json:"name"`
		Description        string    `json:"description"`
		LocationID         int       `json:"location_id"`
		Fqdn               string    `json:"fqdn"`
		Scheme             string    `json:"scheme"`
		BehindProxy        bool      `json:"behind_proxy"`
		MaintenanceMode    bool      `json:"maintenance_mode"`
		Memory             int       `json:"memory"`
		MemoryOverallocate int       `json:"memory_overallocate"`
		Disk               int       `json:"disk"`
		DiskOverallocate   int       `json:"disk_overallocate"`
		UploadSize         int       `json:"upload_size"`
		DaemonListen       int       `json:"daemon_listen"`
		DaemonSftp         int       `json:"daemon_sftp"`
		DaemonBase         string    `json:"daemon_base"`
		CreatedAt          time.Time `json:"created_at"`
		UpdatedAt          time.Time `json:"updated_at"`
	} `json:"attributes"`
}

// NodeEdit is the struct for creating/editing a node for the panel.
type NodeEdit struct {
	Public             bool   `json:"public"`
	Name               string `json:"name"`
	Description        string `json:"description"`
	LocationID         int    `json:"location_id"`
	Fqdn               string `json:"fqdn"`
	Scheme             string `json:"scheme"`
	BehindProxy        bool   `json:"behind_proxy"`
	MaintenanceMode    bool   `json:"maintenance_mode"`
	Memory             int    `json:"memory"`
	MemoryOverallocate int    `json:"memory_overallocate"`
	Disk               int    `json:"disk"`
	DiskOverallocate   int    `json:"disk_overallocate"`
	UploadSize         int    `json:"upload_size"`
	DaemonListen       int    `json:"daemon_listen"`
	DaemonSftp         int    `json:"daemon_sftp"`
	DaemonBase         string `json:"daemon_base"`
}

// NodeAllocations are the allocations for a single node.
type NodeAllocations struct {
	Object string `json:"object"`
	Data   []struct {
		Object     string `json:"object"`
		Attributes struct {
			ID       int    `json:"id"`
			IP       string `json:"ip"`
			Alias    string `json:"alias"`
			Port     int    `json:"port"`
			Assigned bool   `json:"assigned"`
		} `json:"attributes"`
	} `json:"data"`
	Meta struct {
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

// GetNodes returns all available nodes.
func GetNodes() (Nodes, error) {
	var nodes Nodes

	nbytes, err := queryPanelAPI("nodes", "get", nil)
	if err != nil {
		return nodes, err
	}

	// Get node info from the panel
	// Unmarshal the bytes to a usable struct.
	err = json.Unmarshal(nbytes, &nodes)
	if err != nil {
		return nodes, err
	}

	return nodes, nil
}
