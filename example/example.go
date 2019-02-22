package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	croc "github.com/parkervcp/crocgodyl"
)

var (
	// Config is the global config variable.
	Config config
)

type config struct {
	PanelURL    string `json:"panel_url"`
	APIToken    string `json:"api_token"`
	ClientToken string `json:"client_token"`
}

func init() {
	//log.SetOutput(os.Stdout)
	// Open our jsonFile
	jsonFile, err := os.Open("config.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		log.Fatalf("Error loading config.")
	}

	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal(byteValue, &Config)
}

func main() {
	croc.New(Config.PanelURL, Config.ClientToken, Config.APIToken)

	data, err := croc.GetServers()
	if err != nil {
		log.Println("There was an error getting the servers.")
	}

	servers := data.Server

	fmt.Println("All servers on the panel")
	for _, server := range servers {
		fmt.Printf("ID: %d Name: %s\n", server.Attributes.ID, server.Attributes.Name)
	}

	newServerEnvironment := make(map[string]string)

	newServerEnvironment["SERVER_JARFILE"] = "server.jar"
	newServerEnvironment["VANILLA_VERSION"] = "latest"

	newServer := croc.ServerChange{
		Name:        "A Minecraft Server",
		User:        1,
		Egg:         5,
		DockerImage: `quay.io\/pterodactyl\/core:java`,
		Startup:     "java -Xms128M -Xmx {{SERVER_MEMORY}}M -jar {{SERVER_JARFILE}}",
		Environment: newServerEnvironment,
		Limits: croc.ServerLimits{
			Memory: 1024,
			Swap:   0,
			Disk:   1024,
			Io:     500,
			CPU:    0,
		},
		FeatureLimits: croc.ServerFeatureLimits{
			Databases:   0,
			Allocations: 0,
		},
		Allocation: croc.ServerAllocation{
			Default: 2,
		},
	}

	newServerInfo, err := croc.CreateServer(newServer)
	if err != nil {
		log.Println(err)
		os.Exit(37)
	}

	fmt.Println(newServerInfo)
}
