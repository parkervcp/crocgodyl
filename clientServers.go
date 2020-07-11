package crocgodyl

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// Client Servers API

// ClientServers is the default all servers view for the client API.
// GET this from the '/api/client' endpoint
type ClientServers struct {
	Object        string         `json:"object,omitempty"`
	ClientServers []ClientServer `json:"data,omitempty"`
	Meta          Meta           `json:"meta,omitempty"`
}

// ClientServer is the server object view returning single server information.
// GET this from the '/api/client/servers/<server_ID>' endpoint
type ClientServer struct {
	Object     string `json:"object,omitempty"`
	Attributes struct {
		ServerOwner   bool                      `json:"server_owner,omitempty"`
		Identifier    string                    `json:"identifier,omitempty"`
		UUID          string                    `json:"uuid,omitempty"`
		Name          string                    `json:"name,omitempty"`
		Description   string                    `json:"description,omitempty"`
		Limits        ClientServerLimits        `json:"limits,omitempty"`
		FeatureLimits ClientServerFeatureLimits `json:"feature_limits,omitempty"`
	} `json:"attributes,omitempty"`
}

type ClientServerLimits struct {
	Memory int `json:"memory,omitempty"`
	Swap   int `json:"swap,omitempty"`
	Disk   int `json:"disk,omitempty"`
	Io     int `json:"io,omitempty"`
	CPU    int `json:"cpu,omitempty"`
}

type ClientServerFeatureLimits struct {
	Databases   int `json:"databases,omitempty"`
	Allocations int `json:"allocations,omitempty"`
}

// ClientServerUtilization is the server statistics reported by the daemon.
// GET this from the '/api/client/servers/<server_ID>/utilization' endpoint
type ClientServerUtil struct {
	Object     string                 `json:"object,omitempty"`
	Attributes ClientServerUtilAttrib `json:"attributes,omitempty"`
}

type ClientServerUtilAttrib struct {
	State  string                 `json:"state,omitempty"`
	Memory ClientServerUtilLimits `json:"memory,omitempty"`
	CPU    ClientServerUtilLimits `json:"cpu,omitempty"`
	Disk   ClientServerUtilLimits `json:"disk,omitempty"`
}

type ClientServerUtilLimits struct {
	Current float64   `json:"current,omitempty"`
	Cores   []float64 `json:"cores,omitempty"`
	Limit   int       `json:"limit,omitempty"`
}

func (l *ClientServerUtilLimits) UnmarshalJSON(b []byte) error {
	var serverLimits struct {
		Current float64   `json:"current,omitempty"`
		Cores   []float64 `json:"cores,omitempty"`
		Limit   int       `json:"limit,omitempty"`
	}

	if err := json.Unmarshal([]byte(strings.Replace(string(b), "\"cores\":{}", "\"cores\":[]", -1)), &serverLimits); err != nil {
		return err
	}

	l.Current = serverLimits.Current
	l.Cores = serverLimits.Cores
	l.Limit = serverLimits.Limit

	return nil
}

// ClientServerConsoleCommand is the struct for sending a command for the server console
// POST this to the '/api/client/servers/<server_ID>/command' endpoint
type ClientServerConsole struct {
	Command string `json:"command,omitempty"`
}

// ClientServerPowerAction is the struct for sending a power command for the server
// POST this to the '/api/client/servers/<server_ID>/power' endpoint
type ClientServerPowerAction struct {
	Signal string `json:"signal,omitempty"`
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

func (config *ClientConfig) GetClientServerUtilization(serverID string) (serverUtil ClientServerUtil, err error) {
	// Get server info from the panel
	serverBytes, err := config.queryClientAPI(fmt.Sprintf("servers/%s/utilization", serverID), "get", nil)
	if err != nil {
		return
	}

	// Unmarshal the bytes to a usable struct.
	err = json.Unmarshal(serverBytes, &serverUtil)
	if err != nil {
		return
	}

	return
}

func (config *ClientConfig) SendServerCommand(serverID string, command ClientServerConsole) (err error) {
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
