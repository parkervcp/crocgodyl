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

	fmt.Printf("%v\n", user)

	data := user.UpdateDescriptor()
	data.RootAdmin = true
	user, err = app.UpdateUser(user.ID, *data)
	if err != nil {
		handleError(err)
		return
	}

	fmt.Printf("root admin: %v\n\n", user.RootAdmin)

	users, err := app.GetUsers()
	if err != nil {
		handleError(err)
		return
	}

	for i, u := range users {
		fmt.Printf("%d: %v\n", i, u)
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
