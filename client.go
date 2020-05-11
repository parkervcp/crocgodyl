package crocgodyl

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// --------------------------------------------------------------
// Client API

const (
	ServerSignalStart  = "start"
	ServerSignalStop    = "stop"
	ServerSignalRestart = "restart"
	ServerSignalKill    = "kill"
)

// Client Server API

// ClientServers is the default all servers view for the client API.
// GET this from the '/api/client' endpoint
type ClientServers struct {
	Object       string         `json:"object"`
	ClientServer []ClientServer `json:"data"`
	Meta         struct {
		Pagination struct {
			Total       int               `json:"total"`
			Count       int               `json:"count"`
			PerPage     int               `json:"per_page"`
			CurrentPage int               `json:"current_page"`
			TotalPages  int               `json:"total_pages"`
			Links       map[string]string `json:"links"`
		} `json:"pagination"`
	} `json:"meta"`
}

// ClientServer is the server object view returning single server information.
// GET this from the '/api/client/servers/<server_ID>' endpoint
type ClientServer struct {
	Object     string `json:"object"`
	Attributes struct {
		ServerOwner bool   `json:"server_owner"`
		Identifier  string `json:"identifier"`
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
			Databases int `json:"databases"`
		} `json:"feature_limits"`
		Relationships struct {
			Allocations struct {
				Object string `json:"object"`
				Data   []struct {
					Object     string `json:"object"`
					Attributes struct {
						Primary bool   `json:"primary"`
						IP      string `json:"ip"`
						Alias   string `json:"alias"`
						Port    int    `json:"port"`
					} `json:"attributes"`
				} `json:"data"`
			} `json:"allocations"`
		} `json:"relationships"`
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
// GET this from the '/api/client/servers/<server_ID>/command' endpoint
type ClientServerConsoleCommand struct {
	Command string `json:"command"`
}

// ClientServerPowerAction is the struct for sending a power command for the server
// GET this from the '/api/client/servers/<server_ID>/power' endpoint
type ClientServerPowerAction struct {
	Signal string `json:"signal"`
}

// GetClientServers retrieves the servers associated with the client.
func (config *CrocConfig) GetClientServers(pagination bool) (*ClientServers, error) {
	var servers *ClientServers

	// get json bytes from the panel.
	bytes, err := config.queryPanelClient("", "get", nil)
	if err != nil {
		return nil, err
	}

	// Get server info from the panel
	// Unmarshal the bytes to a usable struct.
	err = json.Unmarshal(bytes, &servers)
	if err != nil {
		return nil, err
	}

	if !pagination {
		// Get all pages
		for _, next := range servers.Meta.Pagination.Links {
			res, err := queryURL(next, "get", config.ClientToken, nil)
			if err != nil {
				return nil, err
			}

			bodyBytes, _ := ioutil.ReadAll(res.Body)

			var serverPage *ClientServers
			err = json.Unmarshal(bodyBytes, &serverPage)
			if err != nil {
				return nil, err
			}

			servers.ClientServer = append(servers.ClientServer, serverPage.ClientServer...)
		}
	}


	return servers, nil
}

// GetClientServer retrieves a server associated with the provided id.
func (config *CrocConfig) GetClientServer(serverId string) (*ClientServer, error) {
	var server *ClientServer

	endpoint := fmt.Sprintf("servers/%s?include=allocations", serverId)

	// get json bytes from the panel.
	bytes, err := config.queryPanelClient(endpoint, "get", nil)
	if err != nil {
		return nil, err
	}

	// Get server info from the panel
	// Unmarshal the bytes to a usable struct.
	err = json.Unmarshal(bytes, &server)
	if err != nil {
		return nil, err
	}

	return server, nil
}


// SetServerPowerState changes the power state of a server.
func (config *CrocConfig) SetServerPowerState(serverId string, signal string) error {

	endpoint := fmt.Sprintf("servers/%s/power", serverId)

	svPowerAction := ClientServerPowerAction{Signal: signal}
	body, err := json.Marshal(svPowerAction)
	if err != nil {
		return err
	}

	// get json bytes from the panel.
	_, err = config.queryPanelClient(endpoint, "post", body)
	if err != nil {
		return err
	}

	return nil
}

// SendServerCommand sends a console command to the server of the specified identifier.
func (config *CrocConfig) SendServerCommand(serverId string, command string) error {

	endpoint := fmt.Sprintf("servers/%s/command", serverId)

	svCommand := ClientServerConsoleCommand{Command: command}
	body, err := json.Marshal(svCommand)
	if err != nil {
		return err
	}

	// get json bytes from the panel
	_, err = config.queryPanelClient(endpoint, "post", body)
	if err != nil {
		return err
	}

	return nil
}

// GetClientServerUtilization fetches the usage data for a server.
func (config *CrocConfig) GetClientServerUtilization(serverId string) (*ClientServerUtilization, error) {

	endpoint := fmt.Sprintf("servers/%s/utilization", serverId)

	// get json bytes from the panel
	bytes, err := config.queryPanelClient(endpoint, "get", nil)
	if err != nil {
		return nil, err
	}

	utilization := &ClientServerUtilization{}

	// Get utilization info from the panel
	// Unmarshal the bytes to a usable struct.
	err = json.Unmarshal(bytes, &utilization)
	if err != nil {
		return nil, err
	}

	return utilization, nil
}
