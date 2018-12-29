package crocgodyl

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// VERSION of crocgodyl follows Semantic Versioning. (http://semver.org/)
const VERSION = "0.0.1-alpha"

// --------------------------------------------------------------
// Client API

// Client Server API

// ClientServers is the default all servers view for the client API.
// GET this from the '/api/client' endpoint
type ClientServers struct {
	Object string `json:"object"`
	Data   []struct {
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
// GET this from the '/api/client/servers/<server_ID>/command' endpoint
type ClientServerConsoleCommand struct {
	Command string `json:"command"`
}

// ClientServerPowerAction is the struct for sending a power command for the server
// GET this from the '/api/client/servers/<server_ID>/power' endpoint
type ClientServerPowerAction struct {
	Signal string `json:"signal"`
}

// --------------------------------------------------------------
// Application API

// Application User API

// PanelUsers is the struct for all the panel users.
// GET this from the '/api/application/users` endpoint
type PanelUsers struct {
	Object string `json:"object"`
	Data   []struct {
		Object     string `json:"object"`
		Attributes struct {
			ID         int         `json:"id"`
			ExternalID interface{} `json:"external_id"`
			UUID       string      `json:"uuid"`
			Username   string      `json:"username"`
			Email      string      `json:"email"`
			FirstName  string      `json:"first_name"`
			LastName   string      `json:"last_name"`
			Language   string      `json:"language"`
			RootAdmin  bool        `json:"root_admin"`
			TwoFa      bool        `json:"2fa"`
			CreatedAt  time.Time   `json:"created_at"`
			UpdatedAt  time.Time   `json:"updated_at"`
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

// PanelUser is the struct for all the panel users.
// GET this from the '/api/application/users/<user_ID>` endpoint
type PanelUser struct {
	Object     string `json:"object"`
	Attributes struct {
		ID         int         `json:"id"`
		ExternalID interface{} `json:"external_id"`
		UUID       string      `json:"uuid"`
		Username   string      `json:"username"`
		Email      string      `json:"email"`
		FirstName  string      `json:"first_name"`
		LastName   string      `json:"last_name"`
		Language   string      `json:"language"`
		RootAdmin  bool        `json:"root_admin"`
		TwoFa      bool        `json:"2fa"`
		CreatedAt  time.Time   `json:"created_at"`
		UpdatedAt  time.Time   `json:"updated_at"`
	} `json:"attributes"`
}

// PanelUserEdit is the struct for creating a panel user.
// POST this to the '/api/application/users/` endpoint
// PATCH this to the '/api/application/users/<userID>` endpoint
type PanelUserEdit struct {
	Username   string `json:"username"`
	Email      string `json:"email"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Language   string `json:"language"`
	ExternalID string `json:"external_id"`
	RootAdmin  bool   `json:"root_admin"`
}

// PanelUserCreateResponse is the struct for the response when creating a user.
// POST this from the '/api/application/users/<user_ID>` endpoint
type PanelUserCreateResponse struct {
	Object     string `json:"object"`
	Attributes struct {
		ID         int         `json:"id"`
		ExternalID interface{} `json:"external_id"`
		UUID       string      `json:"uuid"`
		Username   string      `json:"username"`
		Email      string      `json:"email"`
		FirstName  string      `json:"first_name"`
		LastName   string      `json:"last_name"`
		Language   string      `json:"language"`
		RootAdmin  bool        `json:"root_admin"`
		TwoFa      bool        `json:"2fa"`
		CreatedAt  time.Time   `json:"created_at"`
		UpdatedAt  time.Time   `json:"updated_at"`
	} `json:"attributes"`
	Meta struct {
		Resource string `json:"resource"`
	} `json:"meta"`
}

// PanelUserUpdateResponse is the struct for the response when editing a user.
// PATCH this from the '/api/application/users/<user_ID>` endpoint
type PanelUserUpdateResponse struct {
	Object     string `json:"object"`
	Attributes struct {
		ID         int       `json:"id"`
		ExternalID string    `json:"external_id"`
		UUID       string    `json:"uuid"`
		Username   string    `json:"username"`
		Email      string    `json:"email"`
		FirstName  string    `json:"first_name"`
		LastName   string    `json:"last_name"`
		Language   string    `json:"language"`
		RootAdmin  bool      `json:"root_admin"`
		TwoFa      bool      `json:"2fa"`
		CreatedAt  time.Time `json:"created_at"`
		UpdatedAt  time.Time `json:"updated_at"`
	} `json:"attributes"`
}

// Application Location API

// Locations is the struct for all the nodes added to the panel.
// GET this from the '/api/application/locations` endpoint
type Locations struct {
	Object string `json:"object"`
	Data   []struct {
		Object     string `json:"object"`
		Attributes struct {
			ID        int       `json:"id"`
			Short     string    `json:"short"`
			Long      string    `json:"long"`
			UpdatedAt time.Time `json:"updated_at"`
			CreatedAt time.Time `json:"created_at"`
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

// Location is the struct for a single location.
// GET this from the '/api/application/locations/<location_ID>` endpoint
type Location struct {
	Object     string `json:"object"`
	Attributes struct {
		ID        int       `json:"id"`
		Short     string    `json:"short"`
		Long      string    `json:"long"`
		UpdatedAt time.Time `json:"updated_at"`
		CreatedAt time.Time `json:"created_at"`
	} `json:"attributes"`
}

// LocationEdit is the struct for the json when creating a location.
// GET this from the '/api/application/locations/<location_ID>` endpoint
type LocationEdit struct {
	Short string `json:"short"`
	Long  string `json:"long"`
}

// LocatioCreateResponse is the struct for the response when creating a location.
// GET this from the '/api/application/locations/<location_ID>` endpoint
type LocatioCreateResponse struct {
	Object     string `json:"object"`
	Attributes struct {
		ID         int       `json:"id"`
		ExternalID string    `json:"external_id"`
		UUID       string    `json:"uuid"`
		Username   string    `json:"username"`
		Email      string    `json:"email"`
		FirstName  string    `json:"first_name"`
		LastName   string    `json:"last_name"`
		Language   string    `json:"language"`
		RootAdmin  bool      `json:"root_admin"`
		TwoFa      bool      `json:"2fa"`
		CreatedAt  time.Time `json:"created_at"`
		UpdatedAt  time.Time `json:"updated_at"`
	} `json:"attributes"`
}

// Application Node API

// Nodes is the struct for all the nodes added to the panel.
type Nodes struct {
	Object string `json:"object"`
	Data   []struct {
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

// NodeEdit is the struct for creating a node for the panel.
// POST this to the '/api/application/users/` endpoint
// PATCH this to the '/api/application/users/<userID>` endpoint
type NodeEdit struct {
	Name               string `json:"name"`
	LocationID         string `json:"location_id"`
	Fqdn               string `json:"fqdn"`
	Scheme             string `json:"scheme"`
	Memory             int    `json:"memory"`
	MemoryOverallocate int    `json:"memory_overallocate"`
	Disk               int    `json:"disk"`
	DiskOverallocate   int    `json:"disk_overallocate"`
	DaemonListen       string `json:"daemon_listen"`
	DaemonSftp         string `json:"daemon_sftp"`
	DaemonBase         string `json:"daemon_base"`
	BehindProxy        string `json:"behind_proxy"`
	Public             string `json:"public"`
	Throttle           struct {
		Enabled bool `json:"enabled"`
	} `json:"throttle"`
}

var config crocConfig

// crocgodyl structs

type crocConfig struct {
	PanelURL    string
	ClientToken string
	AppToken    string
}

// New sets up the API interface with
func New(panelURL string, clientToken string, appToken string) error {
	if panelURL == "" {
		return errors.New("A panel URL is required to use the API")
	}

	if clientToken == "" && appToken == "" {
		return errors.New("At least one api token is required")
	}

	config.PanelURL = panelURL
	config.ClientToken = clientToken
	config.AppToken = appToken

	return nil
}

func queryPanel(url string) []byte {
	//var for response body byte
	var bodyBytes []byte
	//http get json request
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	//Sets request header for the http request
	req.Header.Add("Authorization", "Bearer "+"")
	//send request
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	if resp.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(resp.Body)
	}
	//set bodyBytes to the response body
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	//Close response thread
	defer resp.Body.Close()
	//return byte structure
	return bodyBytes
}

func printJSON(b []byte) ([]byte, error) {
	var out bytes.Buffer
	err := json.Indent(&out, b, "", "  ")
	return out.Bytes(), err
}

// GetServers returns all available servers.
func GetServers() ([]string, error) {
	var servers []string
	//Print response formated in json
	b, _ := printJSON(queryPanel(config.PanelURL + "/api/application/servers"))
	fmt.Printf("%s", b)

	return servers, nil
}
