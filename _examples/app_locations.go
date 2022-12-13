package main

import (
	"fmt"
	"os"

	croc "github.com/parkervcp/crocgodyl"
)

func main() {
	app, _ := croc.NewApp(os.Getenv("CROC_URL"), os.Getenv("CROC_KEY"))

	loc, err := app.CreateLocation("us", "United States")
	if err != nil {
		handleError(err)
		return
	}

	fmt.Printf("ID: %d - Short: %s - Long: %s\n", loc.ID, loc.Short, loc.Long)

	loc, err = app.UpdateLocation(loc.ID, "us", "United States America")
	if err != nil {
		handleError(err)
		return
	}

	fmt.Printf("ID: %d - Short: %s - Long: %s\n", loc.ID, loc.Short, loc.Long)

	locations, err := app.GetLocations()
	if err != nil {
		handleError(err)
		return
	}

	for _, l := range locations {
		fmt.Printf("%d: %s (%s)\n", l.ID, l.Short, l.Long)
	}

	if err = app.DeleteLocation(loc.ID); err != nil {
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
