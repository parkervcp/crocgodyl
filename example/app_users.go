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
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("%v\n", user)

	data := user.UpdateDescriptor()
	data.RootAdmin = true
	user, err = app.UpdateUser(user.ID, *data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("root admin: %v\n\n", user.RootAdmin)

	users, err := app.GetUsers()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for i, u := range users {
		fmt.Printf("%d: %v\n", i, u)
	}

	if err = app.DeleteUser(user.ID); err != nil {
		fmt.Println(err.Error())
	}
}
