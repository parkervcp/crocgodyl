package crocgodyl

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// Application Node API

// Nodes is the struct for all the nodes added to the panel.
type Nodes struct {
	Object string `json:"object,omitempty"`
	Node   []Node `json:"data,omitempty"`
	Meta   Meta   `json:"meta,omitempty"`
}

// Node is the struct for a single node.
type Node struct {
	Object     string         `json:"object,omitempty"`
	Attributes NodeAttributes `json:"attributes,omitempty"`
}

// NodeAttributes is the struct for the attributes of a node
type NodeAttributes struct {
	ID                 int       `json:"id,omitempty"`
	Public             bool      `json:"public,omitempty"`
	Name               string    `json:"name,omitempty"`
	Description        string    `json:"description,omitempty"`
	LocationID         int       `json:"location_id,omitempty"`
	Fqdn               string    `json:"fqdn,omitempty"`
	Scheme             string    `json:"scheme,omitempty"`
	BehindProxy        bool      `json:"behind_proxy,omitempty"`
	MaintenanceMode    bool      `json:"maintenance_mode,omitempty"`
	Memory             int       `json:"memory,omitempty"`
	MemoryOverallocate int       `json:"memory_overallocate"`
	Disk               int       `json:"disk,omitempty"`
	DiskOverallocate   int       `json:"disk_overallocate"`
	UploadSize         int       `json:"upload_size,omitempty"`
	DaemonListen       int       `json:"daemon_listen,omitempty"`
	DaemonSftp         int       `json:"daemon_sftp,omitempty"`
	DaemonBase         string    `json:"daemon_base,omitempty"`
	CreatedAt          time.Time `json:"created_at,omitempty"`
	UpdatedAt          time.Time `json:"updated_at,omitempty"`
}

// NodeAllocations are the allocations for a single node.
type NodeAllocations struct {
	Object      string       `json:"object,omitempty"`
	Allocations []Allocation `json:"data,omitempty"`
	Meta        Meta         `json:"meta,omitempty"`
}

// Allocation is the struct for an allocation on the node
type Allocation struct {
	Object     string               `json:"object,omitempty"`
	Attributes AllocationAttributes `json:"attributes,omitempty"`
}

// AllocationAttributes is the struct for the allocations on the node.
type AllocationAttributes struct {
	ID       int      `json:"id,omitempty"`
	IP       string   `json:"ip,omitempty"`
	Alias    string   `json:"alias,omitempty"`
	Port     int      `json:"port,omitempty"`
	Ports    []string `json:"ports,omitempty"`
	Assigned bool     `json:"assigned,omitempty"`
}

// GetNodes returns all available nodes.
func GetNodes() (Nodes, error) {
	var nodes Nodes
	endpoint := fmt.Sprintf("nodes/")

	nbytes, err := queryPanelAPI(endpoint, "get", nil)
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

// GetNode inforation on a single node.
// nodeID is an int
func GetNode(nodeID int) (Node, error) {
	var node Node
	endpoint := fmt.Sprintf("nodes/%d", nodeID)

	nbytes, err := queryPanelAPI(endpoint, "get", nil)
	if err != nil {
		return node, err
	}

	// Get node info from the panel
	// Unmarshal the bytes to a usable struct.
	err = json.Unmarshal(nbytes, &node)
	if err != nil {
		return node, err
	}

	return node, nil
}

// GetNodeAllocations information on a single node.
// Depending on how man allocations you have this may take a while.
func GetNodeAllocations(nodeID int) (NodeAllocations, error) {
	var allocations NodeAllocations
	var allocationsAll NodeAllocations

	pages, err := GetNodeAllocationsByPage(nodeID, 1)
	if err != nil {
		return allocations, err
	}

	for i := 1; i <= pages.Meta.Pagination.TotalPages; i++ {
		allocations, err := GetNodeAllocationsByPage(nodeID, i)
		if err != nil {
			return allocationsAll, err
		}

		for _, allocation := range allocations.Allocations {
			allocationsAll.Allocations = append(allocationsAll.Allocations, allocation)
		}
	}

	return allocationsAll, nil
}

// GetNodeAllocationsByPage information on a single node by page count.
// nodeID is an int
func GetNodeAllocationsByPage(nodeID int, pageID int) (NodeAllocations, error) {
	var allocations NodeAllocations
	endpoint := fmt.Sprintf("nodes/%d/allocations?page=%d", nodeID, pageID)

	nabytes, err := queryPanelAPI(endpoint, "get", nil)
	if err != nil {
		return allocations, err
	}

	// Get node info from the panel
	// Unmarshal the bytes to a usable struct.
	err = json.Unmarshal(nabytes, &allocations)
	if err != nil {
		return allocations, err
	}

	return allocations, nil
}

// GetNodeAllocationByPort returns the allocation id and assigned status
func GetNodeAllocationByPort(nodeID int, portNum int) (int, bool, error) {
	allocations, err := GetNodeAllocations(nodeID)
	if err != nil {
		return 0, false, err
	}

	for _, allocation := range allocations.Allocations {
		if allocation.Attributes.Port == portNum {
			return allocation.Attributes.ID, allocation.Attributes.Assigned, nil
		}
	}

	return 0, false, errors.New("port not found")
}

// GetNodeAllocationByID returns the allocation id and assigned status
func GetNodeAllocationByID(nodeID int, allocationID int) (int, bool, error) {
	allocations, err := GetNodeAllocations(nodeID)
	if err != nil {
		return 0, false, err
	}

	for _, allocation := range allocations.Allocations {
		if allocation.Attributes.ID == allocationID {
			return allocation.Attributes.Port, allocation.Attributes.Assigned, nil
		}
	}

	return 0, false, errors.New("id not found")
}

// CreateNode creates a user.
func CreateNode(newNode NodeAttributes) (Node, error) {
	var nodeDetails Node
	endpoint := fmt.Sprintf("nodes/")

	nnbytes, err := json.Marshal(newNode)
	if err != nil {
		return nodeDetails, err
	}

	// get json bytes from the panel.
	nbytes, err := queryPanelAPI(endpoint, "post", nnbytes)
	if err != nil {
		return nodeDetails, err
	}

	// Get server info from the panel
	// Unmarshal the bytes to a usable struct.
	err = json.Unmarshal(nbytes, &nodeDetails)
	if err != nil {
		return nodeDetails, err
	}

	return nodeDetails, nil
}

// CreateNodeAllocations creates a user.
// the panel does not response with a repsonse but a 204
func CreateNodeAllocations(newNodeAllocations AllocationAttributes, nodeID int) error {
	endpoint := fmt.Sprintf("nodes/%d/allocations", nodeID)

	nnabytes, err := json.Marshal(newNodeAllocations)
	if err != nil {
		return err
	}

	// get json bytes from the panel.
	_, err = queryPanelAPI(endpoint, "post", nnabytes)
	if err != nil {
		return err
	}

	return nil
}

// EditNode creates a user.
func EditNode(editNode NodeAttributes, nodeID int) (Node, error) {
	var nodeDetails Node
	endpoint := fmt.Sprintf("nodes/%d", nodeID)

	enbytes, err := json.Marshal(editNode)
	if err != nil {
		return nodeDetails, err
	}

	// get json bytes from the panel.
	nbytes, err := queryPanelAPI(endpoint, "patch", enbytes)
	if err != nil {
		return nodeDetails, err
	}

	// Get server info from the panel
	// Unmarshal the bytes to a usable struct.
	err = json.Unmarshal(nbytes, &nodeDetails)
	if err != nil {
		return nodeDetails, err
	}

	return nodeDetails, nil
}
