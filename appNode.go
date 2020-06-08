package crocgodyl

import (
	"encoding/json"
	"fmt"
	"time"
)

// Application AppNodes API

// AppNodes is the struct for all the nodes added to the panel.
type AppNodes struct {
	Object string `json:"object,omitempty"`
	Nodes  []Node `json:"data,omitempty"`
	Meta   Meta   `json:"meta,omitempty"`
}

// AppNodes is the struct for a single node.
type Node struct {
	Object     string         `json:"object,omitempty"`
	Attributes NodeAttributes `json:"attributes,omitempty"`
}

// NodeAttributes is the struct for the attributes of a node
type NodeAttributes struct {
	ID              int       `json:"id,omitempty"`
	Public          bool      `json:"public,omitempty"`
	Name            string    `json:"name,omitempty"`
	Description     string    `json:"description,omitempty"`
	LocationID      int       `json:"location_id,omitempty"`
	FQDN            string    `json:"fqdn,omitempty"`
	Scheme          string    `json:"scheme,omitempty"`
	BehindProxy     bool      `json:"behind_proxy,omitempty"`
	MaintenanceMode bool      `json:"maintenance_mode,omitempty"`
	Memory          int       `json:"memory,omitempty"`
	MemoryOverAlloc int       `json:"memory_overallocate"`
	Disk            int       `json:"disk,omitempty"`
	DiskOverAlloc   int       `json:"disk_overallocate"`
	UploadSize      int       `json:"upload_size,omitempty"`
	DaemonListen    int       `json:"daemon_listen,omitempty"`
	DaemonSftp      int       `json:"daemon_sftp,omitempty"`
	DaemonBase      string    `json:"daemon_base,omitempty"`
	CreatedAt       time.Time `json:"created_at,omitempty"`
	UpdatedAt       time.Time `json:"updated_at,omitempty"`
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

// GetNodesByPage information on a single node.
// nodeID is an int
func (config *AppConfig) getNodesByPage(pageID int) (nodes AppNodes, err error) {
	nodeBytes, err := config.queryApplicationAPI(fmt.Sprintf("nodes?page=%d", pageID), "get", nil)
	if err != nil {
		return nodes, err
	}

	// Get node info from the panel
	// Unmarshal the bytes to a usable struct.
	if err = json.Unmarshal(nodeBytes, &nodes); err != nil {
		return
	}

	return
}

// GetNodes returns all available nodes.
func (config *AppConfig) GetNodes() (nodes AppNodes, err error) {
	// Get node info from the panel
	nodeBytes, err := config.queryApplicationAPI("nodes", "get", nil)
	if err != nil {
		return
	}

	// Unmarshal the bytes to a usable struct.
	err = json.Unmarshal(nodeBytes, &nodes)
	if err != nil {
		return
	}

	if nodes.Meta.Pagination.TotalPages > 1 {
		for i := 1; i >= nodes.Meta.Pagination.TotalPages; i++ {
			pageNodes, err := config.getNodesByPage(i)
			if err != nil {
				return nodes, err
			}
			for _, location := range pageNodes.Nodes {
				nodes.Nodes = append(nodes.Nodes, location)
			}
		}
	}

	return
}

// GetNode information on a single node.
// nodeID is an int
func (config *AppConfig) GetNode(nodeID int) (node Node, err error) {
	endpoint := fmt.Sprintf("nodes/%d", nodeID)

	nodeBytes, err := config.queryApplicationAPI(endpoint, "get", nil)
	if err != nil {
		return
	}

	// Get node info from the panel
	// Unmarshal the bytes to a usable struct.
	err = json.Unmarshal(nodeBytes, &node)
	if err != nil {
		return
	}

	return
}

// GetNodeAllocationsByPage information on a single node by page count.
// nodeID is an int
func (config *AppConfig) getNodeAllocationsByPage(nodeID int, pageID int) (NodeAllocations, error) {
	var allocations NodeAllocations
	endpoint := fmt.Sprintf("nodes/%d/allocations?page=%d", nodeID, pageID)

	nodeAllocBytes, err := config.queryApplicationAPI(endpoint, "get", nil)
	if err != nil {
		return allocations, err
	}

	// Get node info from the panel
	// Unmarshal the bytes to a usable struct.
	err = json.Unmarshal(nodeAllocBytes, &allocations)
	if err != nil {
		return allocations, err
	}
	return allocations, nil
}

// GetNodeAllocations information on a single node.
// Depending on how man allocations you have this may take a while.
func (config *AppConfig) GetNodeAllocations(nodeID int) (allocations NodeAllocations, err error) {
	// Get allocation info from the panel
	allocBytes, err := config.queryApplicationAPI(fmt.Sprintf("nodes/%d/allocations?page=%d", nodeID), "get", nil)
	if err != nil {
		return
	}

	// Unmarshal the bytes to a usable struct.
	err = json.Unmarshal(allocBytes, &allocations)
	if err != nil {
		return
	}

	for i := 1; i >= allocations.Meta.Pagination.TotalPages; i++ {
		allocations, err := config.getNodeAllocationsByPage(nodeID, i)
		if err != nil {
			return allocations, err
		}
		for _, allocation := range allocations.Allocations {
			allocations.Allocations = append(allocations.Allocations, allocation)
		}
	}

	return
}

// GetNodeAllocationByPort returns the allocation id and assigned status
// returns portID 0 if the port is not assigned to the node.
// returns if the port is used or not
func (config *AppConfig) GetNodeAllocationByPort(nodeID int, portNum int) (portID int, used bool, err error) {
	allocations, err := config.GetNodeAllocations(nodeID)
	if err != nil {
		return
	}

	for _, allocation := range allocations.Allocations {
		if allocation.Attributes.Port == portNum {
			return allocation.Attributes.ID, allocation.Attributes.Assigned, nil
		}
	}

	return
}

// GetNodeAllocationByID returns the allocation id and assigned status
// Takes a node ID and Allocation ID
// returns port 0 if the port is not assigned to the node.
// returns if the port is used or not
func (config *AppConfig) GetNodeAllocationByID(nodeID int, allocationID int) (port int, used bool, err error) {
	allocations, err := config.GetNodeAllocations(nodeID)
	if err != nil {
		return
	}

	for _, allocation := range allocations.Allocations {
		if allocation.Attributes.ID == allocationID {
			return allocation.Attributes.Port, allocation.Attributes.Assigned, nil
		}
	}

	return
}

// Node Modifications

// CreateNode creates a Node.
func (config *AppConfig) CreateNode(newNode NodeAttributes) (node Node, err error) {
	endpoint := fmt.Sprintf("nodes/")

	newNodeBytes, err := json.Marshal(newNode)
	if err != nil {
		return
	}

	// get json bytes from the panel.
	nodeBytes, err := config.queryApplicationAPI(endpoint, "post", newNodeBytes)
	if err != nil {
		return
	}

	// Get server info from the panel
	// Unmarshal the bytes to a usable struct.
	err = json.Unmarshal(nodeBytes, &node)
	if err != nil {
		return
	}
	return
}

// EditNode edits a nodes information.
func (config *AppConfig) EditNode(editNode NodeAttributes, nodeID int) (node Node, err error) {
	endpoint := fmt.Sprintf("nodes/%d", nodeID)

	editNodeBytes, err := json.Marshal(editNode)
	if err != nil {
		return
	}

	// get json bytes from the panel.
	nodeBytes, err := config.queryApplicationAPI(endpoint, "patch", editNodeBytes)
	if err != nil {
		return
	}

	// Get server info from the panel
	// Unmarshal the bytes to a usable struct.
	err = json.Unmarshal(nodeBytes, &node)
	if err != nil {
		return
	}
	return
}

// DeleteNode send a delete request to the panel for a node
// the panel responds with a 204 and no data
// Returns any errors from the panel in json format
func (config *AppConfig) DeleteNode(nodeID int) (err error) {
	endpoint := fmt.Sprintf("nodes/%d", nodeID)

	// get json bytes from the panel.
	_, err = config.queryApplicationAPI(endpoint, "delete", nil)
	if err != nil {
		return err
	}
	return
}

// Allocation Modifications

// CreateNodeAllocations adds Node Allocations.
// the panel responds with a 204 and no data
// Returns any errors from the panel in json format.
func (config *AppConfig) CreateNodeAllocations(newNodeAllocations AllocationAttributes, nodeID int) (err error) {
	endpoint := fmt.Sprintf("nodes/%d/allocations", nodeID)

	newNodeAllocBytes, err := json.Marshal(newNodeAllocations)
	if err != nil {
		return
	}

	// get json bytes from the panel.
	_, err = config.queryApplicationAPI(endpoint, "post", newNodeAllocBytes)
	if err != nil {
		return
	}
	return
}

// DeleteNodeAllocation send a delete request to the panel for a node allocation.
// the panel responds with a 204 and no data
// Returns any errors from the panel in json format.
func (config *AppConfig) DeleteNodeAllocation(nodeID int, allocID int) (err error) {
	endpoint := fmt.Sprintf("nodes/%d/allocations/%d", nodeID, allocID)

	// get json bytes from the panel.
	_, err = config.queryApplicationAPI(endpoint, "delete", nil)
	if err != nil {
		return err
	}
	return
}
