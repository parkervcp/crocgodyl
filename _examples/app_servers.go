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
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("%v\n", server)

	data := server.DetailsDescriptor()
	data.ExternalID = "croc"
	server, err = app.UpdateServerDetails(server.ID, *data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("external id: %s\n\n", server.ExternalID)

	servers, err := app.GetServers()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for i, s := range servers {
		fmt.Printf("%d: %d\n", i, s)
	}

	if err = app.DeleteServer(server.ID, false); err != nil {
		fmt.Println(err.Error())
	}
}
