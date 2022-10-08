package main

import (
	"fmt"
	"os"

	croc "github.com/parkervcp/crocgodyl"
)

func main() {
	app, _ := croc.NewApp(os.Getenv("CROC_URL"), os.Getenv("CROC_KEY"))

	server, err := app.CreateServer(croc.CreateServerDescriptor{
		Name:        "crocgodyl server",
		Description: "test server",
		User:        5,
		Egg:         25,
		DockerImage: "quay.io/parkervcp/pterodactyl-images:base_debian",
		Startup:     "./${EXECUTABLE}",
		Environment: map[string]interface{}{
			"GO_PACKAGE": "github.com/parkervcp/crocgodyl",
			"EXECUTABLE": "crocgodyl",
		},
		Limits:        &croc.Limits{1024, 0, 1024, 10, 1, "0", false},
		FeatureLimtis: croc.FeatureLimits{1, 0, 0},
		Deploy:        &croc.DeployDescriptor{[]int{1, 2}, false, []string{}},
	})
	if err != nil {
		handleError(err)
		return
	}

	fmt.Printf("ID: %d - Name: %s - ExternalID: %s\n", server.ID, server.Name, server.ExternalID)

	data := server.DetailsDescriptor()
	data.ExternalID = "croc"
	server, err = app.UpdateServerDetails(server.ID, *data)
	if err != nil {
		handleError(err)
		return
	}

	fmt.Printf("ID: %d - Name: %s - ExternalID: %s\n", server.ID, server.Name, server.ExternalID)

	servers, err := app.GetServers()
	if err != nil {
		handleError(err)
		return
	}

	for _, s := range servers {
		fmt.Printf("%d: %s\n", s.ID, s.Name)
	}

	if err = app.DeleteServer(server.ID, false); err != nil {
		handleError(err)
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
