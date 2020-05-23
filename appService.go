package crocgodyl

import (
	"encoding/json"
	"time"
)

// Application Service API - Includes Nests and Eggs

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

// GetNests returns all available nodes.
func (config *CrocConfig) GetNests() (Nests, error) {
	var nests Nests

	// get json bytes from the panel.
	nestBytes, err := config.queryPanelAPI("nests", "get", nil)
	if err != nil {
		return nests, err
	}

	// Unmarshal the bytes to a usable struct.
	err = json.Unmarshal(nestBytes, &nests)
	if err != nil {
		return nests, err
	}

	return nests, nil
}

// GetEggs returns all available nodes.
func (config *CrocConfig) GetEggs() (Eggs, error) {
	var eggs Eggs

	return eggs, nil
}
