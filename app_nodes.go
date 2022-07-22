package crocgodyl

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type Node struct {
	ID                 int        `json:"id"`
	Name               string     `json:"name"`
	Description        string     `json:"description"`
	LocationID         int        `json:"location_id"`
	Public             bool       `json:"public"`
	FQDN               string     `json:"fqdn"`
	Scheme             string     `json:"scheme"`
	BehindProxy        bool       `json:"behind_proxy"`
	Memory             int64      `json:"memory"`
	MemoryOverallocate int64      `json:"memory_overallocate"`
	Disk               int64      `json:"disk"`
	DiskOverallocate   int64      `json:"disk_overallocate"`
	DaemonBase         string     `json:"daemon_base"`
	DaemonSftp         int32      `json:"daemon_sftp"`
	DaemonListen       int32      `json:"daemon_listen"`
	MaintenanceMode    bool       `json:"maintenance_mode"`
	UploadSize         int64      `json:"upload_size"`
	CreatedAt          *time.Time `json:"created_at"`
	UpdatedAt          *time.Time `json:"updated_at,omitempty"`
}

func (a *Application) GetNodes() ([]*Node, error) {
	req := a.newRequest("GET", "/nodes", nil)
	res, err := a.Http.Do(req)
	if err != nil {
		return nil, err
	}

	buf, err := validate(res)
	if err != nil {
		return nil, err
	}

	var model struct {
		Data []struct {
			Attributes *Node `json:"attributes"`
		} `json:"data"`
	}
	if err = json.Unmarshal(buf, &model); err != nil {
		return nil, err
	}

	nodes := make([]*Node, 0, len(model.Data))
	for _, n := range model.Data {
		nodes = append(nodes, n.Attributes)
	}

	return nodes, nil
}

func (a *Application) GetNode(id int) (*Node, error) {
	req := a.newRequest("GET", fmt.Sprintf("/nodes/%d", id), nil)
	res, err := a.Http.Do(req)
	if err != nil {
		return nil, err
	}

	buf, err := validate(res)
	if err != nil {
		return nil, err
	}

	var model struct {
		Attributes Node `json:"attributes"`
	}
	if err = json.Unmarshal(buf, &model); err != nil {
		return nil, err
	}

	return &model.Attributes, nil
}

type DeployableNodesDescriptor struct {
	Page         int   `json:"page,omitempty"`
	Memory       int64 `json:"memory"`
	Disk         int64 `json:"disk"`
	LocationsIDs []int `json:"location_ids,omitempty"`
}

func (a *Application) GetDeployableNodes(fields DeployableNodesDescriptor) ([]*Node, error) {
	data, _ := json.Marshal(fields)
	body := bytes.Buffer{}
	body.Write(data)

	req := a.newRequest("GET", "/nodes/deployable", &body)
	res, err := a.Http.Do(req)
	if err != nil {
		return nil, err
	}

	buf, err := validate(res)
	if err != nil {
		return nil, err
	}

	var model struct {
		Data []struct {
			Attributes *Node `json:"attributes"`
		} `json:"data"`
	}
	if err = json.Unmarshal(buf, &model); err != nil {
		return nil, err
	}

	nodes := make([]*Node, 0, len(model.Data))
	for _, n := range model.Data {
		nodes = append(nodes, n.Attributes)
	}

	return nodes, nil
}

type NodeConfiguration struct {
	Debug   bool   `json:"debug"`
	UUID    string `json:"uuid"`
	TokenID string `json:"token_id"`
	Token   string `json:"token"`
	API     struct {
		Host string `json:"host"`
		Port int32  `json:"port"`
		SSL  struct {
			Enabled bool   `json:"enabled"`
			Cert    string `json:"cert"`
			Key     string `json:"key"`
		} `json:"ssl"`
		UploadLimit int64 `json:"upload_limit"`
	} `json:"api"`
	System struct {
		Data string `json:"data"`
		SFTP struct {
			BindPort int32 `json:"bind_port"`
		} `json:"sftp"`
	} `json:"system"`
	AllowedMounts []string `json:"allowed_mounts"`
	Remote        string   `json:"remote"`
}

func (a *Application) GetNodeConfiguration(id int) (*NodeConfiguration, error) {
	req := a.newRequest("GET", fmt.Sprintf("/nodes/%d/configuration", id), nil)
	res, err := a.Http.Do(req)
	if err != nil {
		return nil, err
	}

	buf, err := validate(res)
	if err != nil {
		return nil, err
	}

	var model *NodeConfiguration
	if err = json.Unmarshal(buf, &model); err != nil {
		return nil, err
	}

	return model, nil
}

type CreateNodeDescriptor struct {
	Name               string `json:"name"`
	Description        string `json:"description"`
	LocationID         int    `json:"location_id"`
	Public             bool   `json:"public"`
	FQDN               string `json:"fqdn"`
	Scheme             string `json:"scheme"`
	BehindProxy        bool   `json:"behind_proxy"`
	Memory             int64  `json:"memory"`
	MemoryOverallocate int64  `json:"memory_overallocate"`
	Disk               int64  `json:"disk"`
	DiskOverallocate   int64  `json:"disk_overallocate"`
	DaemonBase         string `json:"daemon_base"`
	DaemonSftp         int32  `json:"daemon_sftp"`
	DaemonListen       int32  `json:"daemon_listen"`
	UploadSize         int64  `json:"upload_size"`
}

func (a *Application) CreateNode(fields CreateNodeDescriptor) (*Node, error) {
	data, _ := json.Marshal(fields)
	body := bytes.Buffer{}
	body.Write(data)

	req := a.newRequest("POST", "/nodes", &body)
	res, err := a.Http.Do(req)
	if err != nil {
		return nil, err
	}

	buf, err := validate(res)
	if err != nil {
		return nil, err
	}

	var model struct {
		Attributes Node `json:"attributes"`
	}
	if err = json.Unmarshal(buf, &model); err != nil {
		return nil, err
	}

	return &model.Attributes, nil
}

type UpdateNodeDescriptor struct {
	Name               string `json:"name,omitempty"`
	Description        string `json:"description,omitempty"`
	LocationID         int    `json:"location_id,omitempty"`
	Public             bool   `json:"public,omitempty"`
	FQDN               string `json:"fqdn,omitempty"`
	Scheme             string `json:"scheme,omitempty"`
	BehindProxy        bool   `json:"behind_proxy,omitempty"`
	Memory             int64  `json:"memory,omitempty"`
	MemoryOverallocate int64  `json:"memory_overallocate,omitempty"`
	Disk               int64  `json:"disk,omitempty"`
	DiskOverallocate   int64  `json:"disk_overallocate,omitempty"`
	DaemonBase         string `json:"daemon_base,omitempty"`
	DaemonSftp         int32  `json:"daemon_sftp,omitempty"`
	DaemonListen       int32  `json:"daemon_listen,omitempty"`
	UploadSize         int64  `json:"upload_size,omitempty"`
}

func (a *Application) UpdateNode(id int, fields UpdateNodeDescriptor) (*Node, error) {
	data, _ := json.Marshal(fields)
	if len(data) == 2 {
		return nil, errors.New("no update fields specified")
	}

	body := bytes.Buffer{}
	body.Write(data)

	req := a.newRequest("PATCH", fmt.Sprintf("/nodes/%d", id), &body)
	res, err := a.Http.Do(req)
	if err != nil {
		return nil, err
	}

	buf, err := validate(res)
	if err != nil {
		return nil, err
	}

	var model struct {
		Attributes Node `json:"attributes"`
	}
	if err = json.Unmarshal(buf, &model); err != nil {
		return nil, err
	}

	return &model.Attributes, nil
}

func (a *Application) DeleteNode(id int) error {
	req := a.newRequest("DELETE", fmt.Sprintf("/nodes/%d", id), nil)
	res, err := a.Http.Do(req)
	if err != nil {
		return err
	}

	_, err = validate(res)
	return err
}

type Allocation struct {
	ID       int    `json:"id"`
	IP       string `json:"ip"`
	Alias    string `json:"alias,omitempty"`
	Port     int32  `json:"port"`
	Notes    string `json:"notes,omitempty"`
	Assigned bool   `json:"assigned"`
}

func (a *Application) GetNodeAllocations(node int) ([]*Allocation, error) {
	req := a.newRequest("GET", fmt.Sprintf("/nodes/%d/allocations", node), nil)
	res, err := a.Http.Do(req)
	if err != nil {
		return nil, err
	}

	buf, err := validate(res)
	if err != nil {
		return nil, err
	}

	var model struct {
		Data []struct {
			Attributes *Allocation `json:"attributes"`
		} `json:"data"`
	}
	if err = json.Unmarshal(buf, &model); err != nil {
		return nil, err
	}

	allocs := make([]*Allocation, 0, len(model.Data))
	for _, a := range model.Data {
		allocs = append(allocs, a.Attributes)
	}

	return allocs, nil
}

type CreateAllocationsDescriptor struct {
	IP    string   `json:"ip"`
	Alias string   `json:"alias,omitempty"`
	Ports []string `json:"ports"`
}

func (a *Application) CreateNodeAllocations(node int, fields CreateAllocationsDescriptor) error {
	data, _ := json.Marshal(fields)
	body := bytes.Buffer{}
	body.Write(data)

	req := a.newRequest("POST", fmt.Sprintf("/nodes/%d/allocations", node), &body)
	res, err := a.Http.Do(req)
	if err != nil {
		return err
	}

	_, err = validate(res)
	return err
}

func (a *Application) DeleteNodeAllocation(node, id int) error {
	req := a.newRequest("DELETE", fmt.Sprintf("/nodes/%d/allocations/%d", node, id), nil)
	res, err := a.Http.Do(req)
	if err != nil {
		return err
	}

	_, err = validate(res)
	return err
}
