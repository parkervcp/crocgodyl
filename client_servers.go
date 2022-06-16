package crocgodyl

import (
	"encoding/json"
)

type Limits struct {
	Memory      int64  `json:"memory"`
	Swap        int64  `json:"swap"`
	Disk        int64  `json:"disk"`
	IO          int64  `json:"io"`
	CPU         int64  `json:"cpu"`
	Threads     string `json:"threads"`
	OOMDisabled bool   `json:"oom_disabled"`
}

type FeatureLimits struct {
	Allocations int `json:"allocations"`
	Backups     int `json:"backups"`
	Databases   int `json:"databases"`
}

type ClientServer struct {
	ServerOwner bool   `json:"server_owner"`
	Identifier  string `json:"identifier"`
	UUID        string `json:"uuid"`
	InternalID  int    `json:"internal_id"`
	Name        string `json:"name"`
	Node        string `json:"node"`
	SFTP        struct {
		IP   string `json:"ip"`
		Port int64  `json:"port"`
	} `json:"sftp_details"`
	Description   string        `json:"description"`
	Limits        Limits        `json:"limits"`
	Invocation    string        `json:"invocation"`
	DockerImage   string        `json:"docker_image"`
	EggFeatures   []string      `json:"egg_features"`
	FeatureLimits FeatureLimits `json:"feature_limits"`
	Status        string        `json:"status"`
	is_suspended  bool
	is_installing bool
	Transferring  bool `json:"is_transferring"`
}

func (s *ClientServer) Suspended() bool {
	return s.is_suspended
}

func (s *ClientServer) Installing() bool {
	return s.is_installing
}

func (c *Client) Servers() ([]ClientServer, error) {
	req := c.newRequest("GET", "", nil)
	res, err := c.Http.Do(req)
	if err != nil {
		return nil, err
	}

	buf, err := validate(res)
	if err != nil {
		return nil, err
	}

	var model *struct {
		Data []struct {
			Attributes ClientServer `json:"attributes"`
		} `json:"data"`
	}
	if err = json.Unmarshal(buf, &model); err != nil {
		return nil, err
	}

	servers := make([]ClientServer, 0, len(model.Data))
	for _, s := range model.Data {
		servers = append(servers, s.Attributes)
	}

	return servers, nil
}

func (c *Client) Server(identifier string) (*ClientServer, error) {
	req := c.newRequest("GET", "/servers/"+identifier, nil)
	res, err := c.Http.Do(req)
	if err != nil {
		return nil, err
	}

	buf, err := validate(res)
	if err != nil {
		return nil, err
	}

	var model *struct {
		Attributes ClientServer `json:"attributes"`
	}
	if err = json.Unmarshal(buf, &model); err != nil {
		return nil, err
	}

	return &model.Attributes, err
}
