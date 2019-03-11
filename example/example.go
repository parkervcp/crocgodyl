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

	/*
	    __ _____ ___ _______
	   / // (_-</ -_) __(_-<
	   \_,_/___/\__/_/ /___/
	*/

	/*
		// get configs and print the user id and usernames of all the users on the panel.
		fmt.Println("All users on the panel")

			userData, err := croc.GetUsers()
		if err != nil {
			log.Println("There was an error getting the users.")
		}

		fmt.Println("All users on the panel")
		for _, user := range userData.User {
			fmt.Printf("ID: %d Name: %s\n", user.Attributes.ID, user.Attributes.Username)
		}

		fmt.Println("Creating user")
		newUser := croc.UserAttributes{
			Email:     "user@domain.tld",
			Username:  "user1234",
			FirstName: "Some",
			LastName:  "User",
		}

		newUserInfo, err := croc.CreateUser(newUser)
		if err != nil {
			log.Println("Failed to create user.")
			log.Println(err)
		} else {
			log.Println("User created successfully.")
			fmt.Println("user info")
			fmt.Printf("ID: %d Username: %s\n", newUserInfo.Attributes.ID, newUserInfo.Attributes.Username)
		}

		fmt.Println("Editing user")
		editUser := croc.UserAttributes{
			Email:      "user@domain.tld",
			Username:   "user1234",
			FirstName:  "Someone",
			LastName:   "Userfied",
			Password:   "aP@55word",
			ExternalID: "99",
		}

		editUserInfo, err := croc.EditUser(editUser, 3)
		if err != nil {
			log.Println("User edit failed.")
			log.Println(err)
		} else {
			log.Println("User edited successfully.")
			fmt.Println("user info")
			fmt.Printf("ID: %d Name: %s\n", editUserInfo.Attributes.ID, editUserInfo.Attributes.Username)
		}

		err := croc.DeleteUser(3)
		if err != nil {
			log.Println(err)
		} else {
			log.Println("User deleted succesfully.")
		}

	*/

	/*
		   __              __  _
		  / /__  _______ _/ /_(_)__  ___  ___
		 / / _ \/ __/ _ `/ __/ / _ \/ _ \(_-<
		/_/\___/\__/\_,_/\__/_/\___/_//_/___/
	*/

	/*
		fmt.Println("Listing all locations on the panel.")
		locationsData, err := croc.GetLocations()
		if err != nil {
			log.Println("There was an error getting the locations.")
		}

		fmt.Println("All users on the panel")
		for _, location := range locationsData.Location {
			fmt.Printf("ID: %d Name: %s\n", location.Attributes.ID, location.Attributes.Long)
		}

		fmt.Println("Listing info on location 1 from the panel")
		locationData, err := croc.GetLocation(1)
		if err != nil {
			log.Println("There was an error getting the locations.")
		}

		fmt.Println("All users on the panel")
		fmt.Printf("ID: %d Name: %s\n", locationData.Attributes.ID, locationData.Attributes.Long)

		newLocation := croc.LocationAttributes{
			Short: "us",
			Long:  "us datacenter",
		}

		newLocationInfo, err := croc.CreateLocation(newLocation)
		if err != nil {
			log.Println("Failed to create location.")
			log.Println(err)
		} else {
			log.Println("location created successfully.")
			fmt.Println("location info")
			fmt.Printf("ID: %d ShortName: %s\n", newLocationInfo.Attributes.ID, newLocationInfo.Attributes.Short)
		}

		editLocation := croc.LocationAttributes{
			Short: "us-la",
			Long:  "us los angelos datacenter",
		}

		editLocationInfo, err := croc.EditLocation(editLocation, 5)
		if err != nil {
			log.Println("Failed to edit location.")
			log.Println(err)
		} else {
			log.Println("location edited successfully.")
			fmt.Println("location info")
			fmt.Printf("ID: %d ShortName: %s\n", editLocationInfo.Attributes.ID, editLocationInfo.Attributes.Short)
		}

		err = croc.DeleteLocation(5)
		if err != nil {
			log.Println(err)
		} else {
			log.Println("location deleted succesfully.")
		}
	*/

	/*
	                   __
	     ___  ___  ___/ /__ ___
	    / _ \/ _ \/ _  / -_|_-<
	   /_//_/\___/\_,_/\__/___/
	*/

<<<<<<< Updated upstream
	// All Nodes
	fmt.Println("Listing all nodes on the panel.")
	nodesData, err := croc.GetNodes()
	if err != nil {
		log.Println("There was an error getting the allocations.")
	}
=======
	/*
		// All Nodes
		fmt.Println("Listing all nodes on the panel.")
		nodesData, err := croc.GetNodes()
		if err != nil {
			log.Println("There was an error getting the locations.")
		}
>>>>>>> Stashed changes

		fmt.Println("All nodes on the panel")
		for _, node := range nodesData.Node {
			fmt.Printf("ID: %d Name: %s\n", node.Attributes.ID, node.Attributes.Name)
		}

		// Single Node
		fmt.Println("Information on Node 1")
		nodeData, err := croc.GetNode(1)
		if err != nil {
			log.Println("There was an error getting the locations.")
		}

		fmt.Printf("ID: %d Name: %s\n", nodeData.Attributes.ID, nodeData.Attributes.Name)

		fmt.Println("Getting allocation id by looking up the port.")
		allocationID, allocationAssigned, err := croc.GetNodeAllocationByPort(2, 25566)
		if err != nil {
			log.Println(err)
		} else {
			fmt.Printf("Allocation id: %d\nAssigned: %t\n", allocationID, allocationAssigned)
		}
		fmt.Println("Getting port by looking up the allocation id.")
		allocationPort, allocationAssigned, err := croc.GetNodeAllocationByID(2, 2)
		if err != nil {
			log.Println(err)
		} else {
			fmt.Printf("Allocation id: %d\nAssigned: %t\n", allocationPort, allocationAssigned)
		}

<<<<<<< Updated upstream
	// Single Node All Ports and allocations
	fmt.Println("Allocations on Node 1")
	nodeAllocData, err := croc.GetNodeAllocations(2)
	if err != nil {
		log.Println(err)
	}

	for _, alloc := range nodeAllocData.Allocations {
		fmt.Printf("ID: %d Port: %d\n", alloc.Attributes.ID, alloc.Attributes.Port)
	}

	// get allocation_id from the port
	fmt.Println("Getting allocation id on node2 by looking up port 25566.")
	allocationID, allocationAssigned, err := croc.GetNodeAllocationByPort(2, 25566)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Printf("Allocation id: %d\nAssigned: %t\n", allocationID, allocationAssigned)
	}

	// get port from the allocation number
	fmt.Println("Getting port on node 2 by looking up allocation id 2.")
	allocationPort, allocationAssigned, err := croc.GetNodeAllocationByID(2, 2)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Printf("Allocation id: %d\nAssigned: %t\n", allocationPort, allocationAssigned)
	}
=======
>>>>>>> Stashed changes

		newNode := croc.NodeAttributes{
			Name:               "testing4",
			LocationID:         3,
			Fqdn:               "testing4.synahost.com",
			Scheme:             "https",
			Memory:             1024,
			MemoryOverallocate: 0,
			Disk:               1024,
			DiskOverallocate:   0,
			DaemonBase:         "/srv/daemon-data",
			DaemonListen:       8080,
			DaemonSftp:         2022,
		}
		newNodeInfo, err := croc.CreateNode(newNode)
		if err != nil {
			log.Println("Failed to create node.")
			log.Println(err)
		} else {
			log.Println("node created successfully.")
			fmt.Println("node info")
			fmt.Printf("ID: %d Name: %s\n", newNodeInfo.Attributes.ID, newNodeInfo.Attributes.Name)
		}
		newNodeAllocations := croc.AllocationAttributes{
			IP:    "2.2.2.2",
			Alias: "two.two.two.two",
			Ports: []string{"4000", "4001", "4002-4500"},
		}

		err = croc.CreateNodeAllocations(newNodeAllocations, 7)
		if err != nil {
			log.Println("Failed to add node allocations.")
			log.Println(err)
		} else {
			log.Println("node allocations added successfully.")
		}

		editNode := croc.NodeAttributes{
			Name:               "testing2",
			LocationID:         3,
			Fqdn:               "testing2.synahost.com",
			Scheme:             "https",
			Memory:             1024,
			MemoryOverallocate: 0,
			Disk:               1024,
			DiskOverallocate:   0,
			DaemonBase:         "/srv/daemon-data",
			DaemonListen:       8080,
			DaemonSftp:         2022,
		}

		editNodeInfo, err := croc.EditNode(editNode, 2)
		if err != nil {
			log.Println("Failed to edit node.")
			log.Println(err)
		} else {
			log.Println("node edited successfully.")
			fmt.Println("node info")
			fmt.Printf("ID: %d Name: %s\n", editNodeInfo.Attributes.ID, editNodeInfo.Attributes.Name)
		}
	*/
	/*
		  ___ ___ _____  _____ _______
		 (_-</ -_) __/ |/ / -_) __(_-<
		/___/\__/_/  |___/\__/_/ /___/
	*/

	// get server information and print the id and names of the first page of servers on the panel.
	fmt.Println("All servers on the panel")

	serversData, err := croc.GetServers()
	if err != nil {
		log.Println("There was an error getting the servers.")
	}

	for _, server := range serversData.Server {
		fmt.Printf("ID: %d Name: %s\n", server.Attributes.ID, server.Attributes.Name)
	}

	// Get information on a single server.
	serverData, err := croc.GetServer(1)
	if err != nil {
		log.Println("There was an error getting the servers.")
	}

	fmt.Printf("ID: %d Name: %s\n", serverData.Attributes.ID, serverData.Attributes.Name)

	var serverPorts []int

	for _, relation := range serverData.Attributes.Relationships.Allocations.Data {
		serverPorts = append(serverPorts, relation.Attributes.Port)
	}

	log.Printf("The server has the following ports assinged: %d\n", serverPorts)

	/*
		// build out a new server config.
		// this is for a vanilla minecraft server.

		fmt.Println("Creating a new server.")
		// The environment variables a map of string keys and values.
		newServerEnvironment := make(map[string]string)

		newServerEnvironment["SERVER_JARFILE"] = "server.jar"
		newServerEnvironment["VANILLA_VERSION"] = "latest"

		// The rest of the server can all be configured as a single struct.
		newServer := croc.ServerChange{
			Name:        "A Minecraft Server",
			User:        1,
			Egg:         5,
			DockerImage: `quay.io/pterodactyl/core:java`,
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

		// Creates a new server and returns the server info that the panel responds with.
		// If there was an error crocgodyle will give you the error code and error message from the panel in json.
		newServerInfo, err := croc.CreateServer(newServer)
		if err != nil {
			log.Println(err)
		} else {
			log.Println("Server created successfully.")
			fmt.Println("New server info")
			fmt.Printf("ID: %d Name: %s\n", newServerInfo.Attributes.ID, newServerInfo.Attributes.Name)
		}

		fmt.Println("Editing server details.")

		editServer := croc.ServerChange{
			Name: "An Awesone Minecraft Server",
			User: 1,
		}

		editedServerInfo, err := croc.EditServerDetails(editServer, 19)
		if err != nil {
			log.Println(err)
		} else {
			fmt.Printf("ID: %d Name: %s\n", editedServerInfo.Attributes.ID, editedServerInfo.Attributes.Name)
		}

		fmt.Println("Editing server build.")

		editServer = croc.ServerChange{
			Name: "An Awesone Minecraft Server",
			User: 1,
		}

		editedServerInfo, err = croc.EditServerDetails(editServer, 19)
		if err != nil {
			log.Println(err)
		} else {
			fmt.Printf("ID: %d Name: %s\n", editedServerInfo.Attributes.ID, editedServerInfo.Attributes.Name)
		}

		// Delete a server
		err = croc.DeleteServer(19)
		if err != nil {
			log.Println(err)
		} else {
			log.Println("Server deleted succesfully.")
		}
	*/
}
