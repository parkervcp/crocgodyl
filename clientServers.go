package crocgodyl

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Client Servers API

// ClientServers is the default all servers view for the client API.
// GET this from the '/api/client' endpoint
type ClientServers struct {
	Object       string         `json:"object"`
	ClientServers []ClientServer `json:"data"`
	Meta         Meta           `json:"meta"`
}

// ClientServer is the server object view returning single server information.
// GET this from the '/api/client/servers/<server_ID>' endpoint
type ClientServer struct {
	Object     string `json:"object"`
	Attributes struct {
		ServerOwner bool   `json:"server_owner"`
		Identifier  string `json:"identifier"`
		UUID        string `json:"uuid"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Limits      struct {
			Memory int `json:"memory"`
			Swap   int `json:"swap"`
			Disk   int `json:"disk"`
			Io     int `json:"io"`
			CPU    int `json:"cpu"`
		} `json:"limits"`
		FeatureLimits struct {
			Databases   int `json:"databases"`
			Allocations int `json:"allocations"`
		} `json:"feature_limits"`
	} `json:"attributes"`
}

// ClientServerUtilization is the server statistics reported by the daemon.
// GET this from the '/api/client/servers/<server_ID>/utilization' endpoint
type ClientServerUtilization struct {
	Object     string `json:"object"`
	Attributes struct {
		State  string `json:"state"`
		Memory struct {
			Current int `json:"current"`
			Limit   int `json:"limit"`
		} `json:"memory"`
		CPU struct {
			Current float64   `json:"current"`
			Cores   []float64 `json:"cores"`
			Limit   int       `json:"limit"`
		} `json:"cpu"`
		Disk struct {
			Current int `json:"current"`
			Limit   int `json:"limit"`
		} `json:"disk"`
	} `json:"attributes"`
}

// ClientServerConsoleCommand is the struct for sending a command for the server console
// POST this to the '/api/client/servers/<server_ID>/command' endpoint
type ClientServerConsoleCommand struct {
	Command string `json:"command"`
}

// ClientServerPowerAction is the struct for sending a power command for the server
// POST this to the '/api/client/servers/<server_ID>/power' endpoint
type ClientServerPowerAction struct {
	Signal string `json:"signal"`
}

func (config *ClientConfig) getClientServersByPage(pageID int) (servers ClientServers, err error) {
	// Get server info from the panel
	serverBytes, err := config.queryClientAPI(fmt.Sprintf("?page=%d", pageID), "get", nil)
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

func (config *ClientConfig) GetClientServers() (servers ClientServers, err error) {
	// Get server info from the panel
	serverBytes, err := config.queryClientAPI("", "get", nil)
	if err != nil {
		return
	}

	// Unmarshal the bytes to a usable struct.
	err = json.Unmarshal(serverBytes, &servers)
	if err != nil {
		return
	}

	if servers.Meta.Pagination.TotalPages > 1 {
		for i := 1; i <= servers.Meta.Pagination.TotalPages; i++ {
			pageServers, err := config.getClientServersByPage(i)
			if err != nil {
				return servers, err
			}
			for _, server := range pageServers.ClientServers {
				servers.ClientServers = append(servers.ClientServers, server)
			}
		}
	}

	return
}

func (config *ClientConfig) GetClientServer(serverID string) (servers ClientServers, err error) {
	// Get server info from the panel
	serverBytes, err := config.queryClientAPI(fmt.Sprintf("servers/%s", serverID), "get", nil)
	if err != nil {
		return
	}

	// Unmarshal the bytes to a usable struct.
	err = json.Unmarshal(serverBytes, &servers)
	if err != nil {
		return
	}

	if servers.Meta.Pagination.TotalPages > 1 {
		for i := 1; i <= servers.Meta.Pagination.TotalPages; i++ {
			pageServers, err := config.getClientServersByPage(i)
			if err != nil {
				return servers, err
			}
			for _, server := range pageServers.ClientServers {
				servers.ClientServers = append(servers.ClientServers, server)
			}
		}
	}

	return
}

func (config *ClientConfig) SendServerCommand(serverID string, command ClientServerConsoleCommand) (err error) {
	commandBytes, err := json.Marshal(command)
	if err != nil {
		return
	}

	// Get server info from the panel
	if _, err = config.queryClientAPI(fmt.Sprintf("servers/%s/command", serverID), "post", commandBytes); err != nil {
		return
	}

	return
}

func (config *ClientConfig) SendServerPowerSignal(serverID, signal string) (err error) {
	powerAction := ClientServerPowerAction{}
	switch signal {
	case "start":
		powerAction.Signal = "start"
	case "stop":
		powerAction.Signal = "stop"
	case "restart":
		powerAction.Signal = "restart"
	case "kill":
		powerAction.Signal = "kill"
	default:
		err = errors.New("power command must be start, stop, restart, or kill")
		return
	}

	powerBytes, err := json.Marshal(powerAction)
	if err != nil {
		return
	}

	// Get server info from the panel
	if _, err = config.queryClientAPI(fmt.Sprintf("servers/%s/power", serverID), "post", powerBytes); err != nil {
		return
	}

	return
}