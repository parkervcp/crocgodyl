package main

import (
	"fmt"
	"os"

	croc "github.com/parkervcp/crocgodyl"
)

func main() {
	app, _ := croc.NewApp(os.Getenv("CROC_URL"), os.Getenv("CROC_KEY"))

	user, err := app.CreateUser(croc.CreateUserDescriptor{
		Email:     "example@example.com",
		Username:  "example",
		FirstName: "test",
		LastName:  "user",
	})
	if err != nil {
		handleError(err)
		return
	}

	fmt.Printf("ID: %d - Name: %s - RootAdmin: %v\n", user.ID, user.Username, user.RootAdmin)

	data := user.UpdateDescriptor()
	data.RootAdmin = true
	user, err = app.UpdateUser(user.ID, *data)
	if err != nil {
		handleError(err)
		return
	}

	fmt.Printf("ID: %d - Name: %s - RootAdmin: %v\n", user.ID, user.Username, user.RootAdmin)

	users, err := app.GetUsers()
	if err != nil {
		handleError(err)
		return
	}

	for _, u := range users {
		fmt.Printf("%d: %s\n", u.ID, u.Username)
	}

	if err = app.DeleteUser(user.ID); err != nil {
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
