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

	for _, server := range servers {
		fmt.Printf("ID: %d Name: %s\n", server.Attributes.ID, server.Attributes.Name)
	}
}
