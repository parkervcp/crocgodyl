package main

import (
	"fmt"
	"os"

	croc "github.com/parkervcp/crocgodyl"
)

func main() {
	client, _ := croc.NewClient(os.Getenv("CROC_URL"), os.Getenv("CROC_KEY"))

	account, err := client.GetAccount()
	if err != nil {
		handleError(err)
		return
	}

	fmt.Printf("ID: %d - Name: %s\n", account.ID, account.FullName())

	apikeys, err := client.GetApiKeys()
	if err != nil {
		handleError(err)
		return
	}

	for _, k := range apikeys {
		fmt.Printf("%s: %s\n", k.Identifier, k.Description)
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
