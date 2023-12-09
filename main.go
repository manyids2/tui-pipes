package main

import (
	"log"

	"github.com/manyids2/tui-pipes/components"
)

func main() {
	app := components.NewApp()
	err := app.App.Run()
	if err != nil {
		log.Fatal(err)
	}
}
