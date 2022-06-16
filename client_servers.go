package crocgodyl

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	Suspended     bool          `json:"is_suspended"`
	Installing    bool          `json:"is_installing"`
	Transferring  bool          `json:"is_transferring"`
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

type WebSocketAuth struct {
	Socket string `json:"socket"`
	Token  string `json:"token"`
}

func (c *Client) ServerWebSocket(identifier string) (*WebSocketAuth, error) {
	req := c.newRequest("GET", fmt.Sprintf("/servers/%s/websocket", identifier), nil)
	res, err := c.Http.Do(req)
	if err != nil {
		return nil, err
	}

	buf, err := validate(res)
	if err != nil {
		return nil, err
	}

	var model *struct {
		Data WebSocketAuth `json:"data"`
	}
	if err = json.Unmarshal(buf, &model); err != nil {
		return nil, err
	}

	return &model.Data, nil
}

type Resources struct {
	MemoryBytes    int64 `json:"memory_bytes"`
	DiskBytes      int64 `json:"disk_bytes"`
	CPUAbsolute    int64 `json:"cpu_absolute"`
	NetworkRxBytes int64 `json:"network_rx_bytes"`
	NetworkTxBytes int64 `json:"network_tx_bytes"`
	Uptime         int64 `json:"uptime"`
}

type Stats struct {
	State     string    `json:"current_state,omitempty"`
	Suspended bool      `json:"is_suspended"`
	Resources Resources `json:"resources"`
}

func (c *Client) ServerStatistics(identifier string) (*Stats, error) {
	req := c.newRequest("GET", fmt.Sprintf("/servers/%s/resources", identifier), nil)
	res, err := c.Http.Do(req)
	if err != nil {
		return nil, err
	}

	buf, err := validate(res)
	if err != nil {
		return nil, err
	}

	var model *struct {
		Attributes Stats `json:"attributes"`
	}
	if err = json.Unmarshal(buf, &model); err != nil {
		return nil, err
	}

	return &model.Attributes, nil
}

func (c *Client) SendServerCommand(identifier, command string) error {
	data, _ := json.Marshal(map[string]string{"command": command})
	body := bytes.Buffer{}
	body.Write(data)

	req := c.newRequest("POST", fmt.Sprintf("/servers/%s/command", identifier), &body)
	res, err := c.Http.Do(req)
	if err != nil {
		return err
	}

	if _, err = validate(res); err != nil {
		return err
	}

	return nil
}

func (c *Client) SetServerPowerState(identifier, state string) error {
	data, _ := json.Marshal(map[string]string{"signal": state})
	body := bytes.Buffer{}
	body.Write(data)

	req := c.newRequest("POST", fmt.Sprintf("/servers/%s/power", identifier), &body)
	res, err := c.Http.Do(req)
	if err != nil {
		return err
	}

	if _, err = validate(res); err != nil {
		return err
	}

	return nil
}
