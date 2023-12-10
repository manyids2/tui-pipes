package main

import (
	"log"

	"github.com/manyids2/tui-pipes/components"
)

func main() {
	config, err := components.ReadConfig("./configs/git_status.json")
	if err != nil {
		log.Fatal(err)
	}

	app := components.NewApp(config)
	err = app.App.Run()
	if err != nil {
		log.Fatal(err)
	}
}
