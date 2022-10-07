package main

import (
	"fmt"
	"os"
	"strings"

	croc "github.com/parkervcp/crocgodyl"
)

func main() {
	url := os.Getenv("CROC_URL")
	app, _ := croc.NewApp(url, os.Getenv("CROC_KEY"))

	node, err := app.CreateNode(croc.CreateNodeDescriptor{
		Name:         "croc-node-1",
		LocationID:   1,
		Public:       true,
		FQDN:         fmt.Sprintf("test.nodes.%s", strings.Split(url, "//")[1]),
		Scheme:       "https",
		BehindProxy:  false,
		Memory:       16000,
		Disk:         1024,
		DaemonBase:   "/var/lib/pterodactyl/volumes",
		DaemonSftp:   2022,
		DaemonListen: 8080,
		UploadSize:   100,
	})
	if err != nil {
		handleError(err)
		return
	}

	fmt.Printf("%v\n", node)

	data := node.UpdateDescriptor()
	data.Public = false
	node, err = app.UpdateNode(node.ID, *data)
	if err != nil {
		handleError(err)
		return
	}

	fmt.Printf("public: %v\n", node.Public)

	nodes, err := app.GetNodes()
	if err != nil {
		handleError(err)
		return
	}

	for i, n := range nodes {
		fmt.Printf("%d: %v\n", i, n)
	}

	if err = app.DeleteNode(node.ID); err != nil {
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
