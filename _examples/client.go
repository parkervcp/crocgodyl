package main

import (
	"fmt"

	croc "github.com/parkervcp/crocgodyl"
)

func client_test() {
	// initialize the client
	client, _ := croc.NewClient("https://pterodactyl.domain", "ptlc_someLongAP1ke3")

	// fetches the client account
	account, err := client.GetAccount()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("%s - %s\n", account.FullName(), account.Email)

	// fetches the servers the account has access to
	servers, err := client.GetServers()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for _, server := range servers {
		fmt.Printf("%s - %s\n", server.Identifier, server.Name)
	}
	server := servers[0]

	// gets the server websocket authentication details for a server
	auth, err := client.GetServerWebSocket(server.Identifier)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("socket: %s\ntoken: %s\n", auth.Socket, auth.Token)

	// gets available api keys associated with the account
	keys, err := client.GetApiKeys()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for _, key := range keys {
		fmt.Printf("%s - %s\n", key.Identifier, key.Description)
	}

	// gets the files in the root directory on a specified server
	files, err := client.GetServerFiles(server.Identifier, "/")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for _, file := range files {
		fmt.Printf("%s - %d bytes\n", file.Name, file.Size)
	}
	file := files[0]

	// gets a downloader object which downloads a file from a server
	dl, err := client.DownloadServerFile(server.Identifier, file.Name)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(dl.URL())

	// set the path on the system to save the download to
	dl.Path = "/" + file.Name

	// execute the downloader and hope it doesn't error
	err = dl.Execute()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
