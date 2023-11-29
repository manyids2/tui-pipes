package main

import (
	"fmt"
	"log"
)

func main() {
	// Get path from args
	path := "layouts/home.yaml"
	fmt.Println(path)
	app := NewApp()

	// Run the application
	err := app.App.Run()
	if err != nil {
		log.Fatal(err)
	}
}
