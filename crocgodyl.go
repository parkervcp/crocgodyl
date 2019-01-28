package crocgodyl

import (
	"bytes"
	"encoding/json"
	"errors"
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
	Object       string         `json:"object"`
	ClientServer []ClientServer `json:"data"`
	Meta         struct {
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

// Users is the struct for all the panel users.
// GET this from the '/api/application/users` endpoint
type Users struct {
	Object string `json:"object"`
	User   []User `json:"data"`
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

// User is the struct for all the panel users.
// GET this from the '/api/application/users/<user_ID>` endpoint
type User struct {
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

// UserCreateResponse is the struct for the response when creating a user.
// POST this from the '/api/application/users/<user_ID>` endpoint
type UserCreateResponse struct {
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

// UserUpdateResponse is the struct for the response when editing a user.
// PATCH this from the '/api/application/users/<user_ID>` endpoint
type UserUpdateResponse struct {
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
	Object   string     `json:"object"`
	Location []Location `json:"data"`
	Meta     struct {
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

// LocationEdit is the struct for the json when editing/creating a location.
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

// Servers is the struct for the servers on the panel.
type Servers struct {
	Object string   `json:"object"`
	Server []Server `json:"data"`
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

// Server is the struct for a server on the panel.
type Server struct {
	Object     string `json:"object"`
	Attributes struct {
		ID          int         `json:"id"`
		ExternalID  interface{} `json:"external_id"`
		UUID        string      `json:"uuid"`
		Identifier  string      `json:"identifier"`
		Name        string      `json:"name"`
		Description string      `json:"description"`
		Suspended   bool        `json:"suspended"`
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
		User       int         `json:"user"`
		Node       int         `json:"node"`
		Allocation int         `json:"allocation"`
		Nest       int         `json:"nest"`
		Egg        int         `json:"egg"`
		Pack       interface{} `json:"pack"`
		Container  struct {
			StartupCommand string                 `json:"startup_command"`
			Image          string                 `json:"image"`
			Installed      bool                   `json:"installed"`
			Environment    map[string]interface{} `json:"environment"`
		} `json:"container"`
		UpdatedAt time.Time `json:"updated_at"`
		CreatedAt time.Time `json:"created_at"`
	} `json:"attributes"`
}

// ServerEdit is the struct for creating/editing a server for the panel.
type ServerEdit struct {
	Name        string   `json:"name"`
	User        int      `json:"user"`
	Egg         int      `json:"egg"`
	DockerImage string   `json:"docker_image"`
	Startup     string   `json:"startup"`
	Environment []string `json:"environment"`
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
	Allocation struct {
		Default int `json:"default"`
	} `json:"allocation"`
}

// Nests is the struct for the nests on the panel.
type Nests struct {
	Object string `json:"object"`
	Nest   []Nest `json:"data"`
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

// Nest is the struct for a nest on the panel.
type Nest struct {
	Object     string `json:"object"`
	Attributes struct {
		ID          int       `json:"id"`
		UUID        string    `json:"uuid"`
		Author      string    `json:"author"`
		Name        string    `json:"name"`
		Description string    `json:"description"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
	} `json:"attributes"`
}

// Eggs is the struct for all eggs in a nest.
type Eggs struct {
	Object string `json:"object"`
	Egg    []Egg  `json:"data"`
}

// Egg is the struct for an egg on the panel.
type Egg struct {
	Object     string `json:"object"`
	Attributes struct {
		ID          int    `json:"id"`
		UUID        string `json:"uuid"`
		Nest        int    `json:"nest"`
		Author      string `json:"author"`
		Description string `json:"description"`
		DockerImage string `json:"docker_image"`
		Config      struct {
			Files struct {
				ConfigYml struct {
					Parser string `json:"parser"`
					Find   struct {
						Listeners0QueryEnabled bool   `json:"listeners[0].query_enabled"`
						Listeners0QueryPort    string `json:"listeners[0].query_port"`
						Listeners0Host         string `json:"listeners[0].host"`
						ServersAddress         struct {
							One27111  string `json:"127.1.1.1"`
							Localhost string `json:"localhost"`
						} `json:"servers.*.address"`
					} `json:"find"`
				} `json:"config.yml"`
			} `json:"files"`
			Startup struct {
				Done            string   `json:"done"`
				UserInteraction []string `json:"userInteraction"`
			} `json:"startup"`
			Stop string `json:"stop"`
			Logs struct {
				Custom   bool   `json:"custom"`
				Location string `json:"location"`
			} `json:"logs"`
			Extends interface{} `json:"extends"`
		} `json:"config"`
		Startup string `json:"startup"`
		Script  struct {
			Privileged bool        `json:"privileged"`
			Install    string      `json:"install"`
			Entry      string      `json:"entry"`
			Container  string      `json:"container"`
			Extends    interface{} `json:"extends"`
		} `json:"script"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	} `json:"attributes"`
}

// ErrorResponse is the response from the panels for errors.
type ErrorResponse struct {
	Errors []struct {
		Code   string `json:"code"`
		Detail string `json:"detail"`
		Source struct {
			Field string `json:"field"`
		} `json:"source"`
	} `json:"errors"`
}

//
// crocgodyl
//
var config crocConfig

type crocConfig struct {
	PanelURL    string
	ClientToken string
	AppToken    string
}

//
// Application code
//

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

func queryPanelAPI(endpoint string) ([]byte, error) {
	//var for response body byte
	var bodyBytes []byte
	//http get json request
	client := &http.Client{}
	req, _ := http.NewRequest("GET", config.PanelURL+"/api/application/"+endpoint, nil)
	//Sets request header for the http request
	req.Header.Add("Authorization", "Bearer "+config.AppToken)
	//send request
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if resp.StatusCode == 422 {
		bodyBytes, _ = ioutil.ReadAll(resp.Body)
		return bodyBytes, errors.New("422")
	}

	if resp.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(resp.Body)
	}
	//set bodyBytes to the response body
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	//Close response thread
	defer resp.Body.Close()

	//return byte structure
	return bodyBytes, nil
}

func printJSON(b []byte) ([]byte, error) {
	var out bytes.Buffer
	err := json.Indent(&out, b, "", "  ")
	return out.Bytes(), err
}

//
// Get Functions
//

// GetLocations returns all available nodes.
func GetLocations() (Locations, error) {
	var locations Locations

	lbytes, err := queryPanelAPI("locations")
	if err != nil {
		return locations, err
	}

	// Get node info from the panel
	// Unmarshal the bytes to a usable struct.
	err = json.Unmarshal(lbytes, &locations)
	if err != nil {
		return locations, err
	}

	return locations, nil
}

// GetNodes returns all available nodes.
func GetNodes() (Nodes, error) {
	var nodes Nodes

	nbytes, err := queryPanelAPI("nodes")
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

// GetServers returns all available servers.
func GetServers() (Servers, error) {
	var servers Servers

	// get json bytes from the panel.
	sbytes, err := queryPanelAPI("servers")
	if err != nil {
		return servers, err
	}

	// Get server info from the panel
	// Unmarshal the bytes to a usable struct.
	err = json.Unmarshal(sbytes, &servers)
	if err != nil {
		return servers, err
	}

	return servers, nil
}

// GetNests returns all available nodes.
func GetNests() (Nests, error) {
	var nests Nests

	// get json bytes from the panel.
	nbytes, err := queryPanelAPI("nests")
	if err != nil {
		return nests, err
	}

	// Unmarshal the bytes to a usable struct.
	err = json.Unmarshal(nbytes, &nests)
	if err != nil {
		return nests, err
	}

	return nests, nil
}

// GetEggs returns all available nodes.
func GetEggs() (Eggs, error) {
	var eggs Eggs

	return eggs, nil
}

// GetUsers returns all available nodes.
func GetUsers() (Users, error) {
	var users Users

	// get json bytes from the panel.
	ubytes, err := queryPanelAPI("users")
	if err != nil {
		return users, err
	}

	// Unmarshal the bytes to a usable struct.
	err = json.Unmarshal(ubytes, &users)
	if err != nil {
		return users, err
	}

	return users, nil
}

//
// Set Functions
//
func SetUser() {

}
