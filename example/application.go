package main

import (
	"fmt"

	croc "github.com/parkervcp/crocgodyl"
)

func main() {
	// initialize the application
	app, _ := croc.NewApp("https://pterodactyl.domain", "ptla_someLongAP1ke3")

	// gets user accounts from the panel
	users, err := app.Users()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for _, user := range users {
		fmt.Println(user.FullName())
	}
	user := users[0]

	// gets server objects from the panel
	servers, err := app.Servers()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for _, server := range servers {
		if server.User == user.ID {
			fmt.Printf("%d - %s\n", server.ID, server.Name)
		}
	}

	// gets a single node from the panel (`app.Nodes()` returns all nodes)
	node, err := app.Node(1)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("%d - %s\n", node.ID, node.Name)

	// returns a single location obhect from the panel (`app.Locations()` returns all locations)
	location, err := app.Location(node.LocationID)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("%d - %s: %s\n", location.ID, location.Short, location.Long)
}
