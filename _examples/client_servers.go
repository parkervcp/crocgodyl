package main

import (
	"fmt"
	"os"

	croc "github.com/parkervcp/crocgodyl"
)

func main() {
	client, _ := croc.NewClient(os.Getenv("CROC_URL"), os.Getenv("CROC_KEY"))

	servers, err := client.GetServers()
	if err != nil {
		handleError(err)
		return
	}

	for _, s := range servers {
		fmt.Printf("%s (%d): %s\n", s.Identifier, s.InternalID, s.Name)
	}

	if len(servers) == 0 {
		fmt.Println("no servers to list")
		return
	}
	server := servers[0]

	if err = client.SetServerPowerState(server.Identifier, "restart"); err != nil {
		handleError(err)
		return
	}

	if err = client.SendServerCommand(server.Identifier, "say \"hello world\""); err != nil {
		handleError(err)
		return
	}
}

func handleError(err error) {
	if errs, ok := err.(*croc.ApiError); ok {
		for _, e := range errs.Errors {
			fmt.Println(e.Error())
		}
	} else {
		fmt.Println(err.Error())
	}
}
