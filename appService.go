package crocgodyl

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Application Service API - Includes Nests and Eggs

// Nests is the struct for the nests on the panel.
type Nests struct {
	Object string `json:"object,omitempty"`
	Nest   []Nest `json:"data,omitempty"`
	Meta   Meta   `json:"meta,omitempty"`
}

// Nest is the struct for a nest on the panel.
type Nest struct {
	Object     string `json:"object,omitempty"`
	Attributes struct {
		ID          int       `json:"id,omitempty"`
		UUID        string    `json:"uuid,omitempty"`
		Author      string    `json:"author,omitempty"`
		Name        string    `json:"name,omitempty"`
		Description string    `json:"description,omitempty"`
		CreatedAt   time.Time `json:"created_at,omitempty"`
		UpdatedAt   time.Time `json:"updated_at,omitempty"`
	} `json:"attributes,omitempty"`
}

// Eggs is the struct for all eggs in a nest.
type NestEggs struct {
	Object string `json:"object,omitempty"`
	Eggs   []Egg  `json:"data,omitempty"`
}

type Egg struct {
	Object     string        `json:"object,omitempty"`
	Attributes EggAttributes `json:"attributes,omitempty"`
}

type EggAttributes struct {
	ID            int          `json:"id,omitempty"`
	UUID          string       `json:"uuid,omitempty"`
	Name          string       `json:"name,omitempty"`
	Nest          int          `json:"nest,omitempty"`
	Author        string       `json:"author,omitempty"`
	Description   string       `json:"description,omitempty"`
	DockerImage   string       `json:"docker_image,omitempty"`
	Config        EggConfig    `json:"config,omitempty"`
	Startup       string       `json:"startup,omitempty"`
	Script        EggScript    `json:"script,omitempty"`
	CreatedAt     time.Time    `json:"created_at,omitempty"`
	UpdatedAt     time.Time    `json:"updated_at,omitempty"`
	Relationships EggRelations `json:"relationships,omitempty"`
}

type EggConfig struct {
	Files   map[string]EggFileConfig `json:"files,omitempty"`
	Startup EggStartup               `json:"startup,omitempty"`
	Stop    string                   `json:"stop,omitempty"`
	Logs    EggLogs                  `json:"logs,omitempty"`
	Extends interface{}              `json:"extends,omitempty"`
}

type EggFileConfig struct {
	Parser string            `json:"parser,omitempty"`
	Find   map[string]string `json:"find,omitempty"`
}

type EggStartup struct {
	Done            string   `json:"done,omitempty"`
	UserInteraction []string `json:"userInteraction,omitempty"`
}

func (e *EggStartup) UnmarshalJSON(b []byte) error {
	var startup struct {
		Done            string   `json:"done,omitempty"`
		UserInteraction []string `json:"userInteraction,omitempty"`
	}

	if err := json.Unmarshal([]byte(strings.Replace(string(b), "\"userInteraction\":{}", "\"userInteraction\":[]", -1)), &startup); err != nil {
		return err
	}

	e.UserInteraction = startup.UserInteraction
	e.Done = startup.Done

	return nil
}

type EggLogs struct {
	Custom   bool   `json:"custom,omitempty"`
	Location string `json:"location,omitempty"`
}

type EggScript struct {
	Privileged bool        `json:"privileged,omitempty"`
	Install    string      `json:"install,omitempty"`
	Entry      string      `json:"entry,omitempty"`
	Container  string      `json:"container,omitempty"`
	Extends    interface{} `json:"extends,omitempty"`
}

type EggRelations struct {
	Variables EggVariables `json:"variables,omitempty"`
}

type EggVariables struct {
	Object string            `json:"object,omitempty"`
	Data   []EggRelationData `json:"data,omitempty"`
}

type EggRelationData struct {
	Object     string `json:"object,omitempty"`
	Attributes struct {
		ID           int    `json:"id,omitempty"`
		EggID        int    `json:"egg_id,omitempty"`
		Name         string `json:"name,omitempty"`
		Description  string `json:"description,omitempty"`
		EnvVariable  string `json:"env_variable,omitempty"`
		DefaultValue string `json:"default_value,omitempty"`
		UserViewable int    `json:"user_viewable,omitempty"`
		UserEditable int    `json:"user_editable,omitempty"`
		Rules        string `json:"rules,omitempty"`
		CreatedAt    string `json:"created_at,omitempty"`
		UpdatedAt    string `json:"updated_at,omitempty"`
	} `json:"attributes,omitempty"`
}

func (e *EggVariables) UnmarshalJSON(b []byte) error {
	var eggVariables struct {
		Object string            `json:"object,omitempty"`
		Data   []EggRelationData `json:"data,omitempty"`
	}

	if err := json.Unmarshal(b, &eggVariables); err != nil {
		if eggVariables.Object == "list" {
			e.Data = []EggRelationData{}
		} else {
			return err
		}
	}

	e.Object = eggVariables.Object
	e.Data = eggVariables.Data

	return nil
}

// GetNests returns all available nodes.
func (config *CrocConfig) GetNests() (nests Nests, err error) {
	// get json bytes from the panel.
	nestBytes, err := config.queryApplicationAPI("nests", "get", nil)
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
func (config *CrocConfig) GetNestEggs(nestID int) (eggs NestEggs, err error) {
	// get json bytes from the panel.
	nestEggsBytes, err := config.queryApplicationAPI(fmt.Sprintf("nests/%d/eggs", nestID), "get", nil)
	if err != nil {
		return eggs, err
	}

	// Unmarshal the bytes to a usable struct.
	err = json.Unmarshal(nestEggsBytes, &eggs)
	if err != nil {
		return eggs, err
	}

	return eggs, nil
}

// GetEggs returns all available nodes.
func (config *CrocConfig) GetEgg(nestID, eggID int) (egg Egg, err error) {
	// get json bytes from the panel.
	eggBytes, err := config.queryApplicationAPI(fmt.Sprintf("nests/%d/eggs/%d?include=variables", nestID, eggID), "get", nil)
	if err != nil {
		return egg, err
	}

	// Unmarshal the bytes to a usable struct.
	if err = json.Unmarshal(eggBytes, &egg); err != nil {
		return
	}

	return egg, nil
}
