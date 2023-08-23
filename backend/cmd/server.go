package main

import (
	"flag"
	"log"

	"github.com/codescalersinternships/Flyspray/app"
)

func main() {
	var configFilePath string
	flag.StringVar(&configFilePath, "f", "config.json", "Specify the filepath of json configuration file")
	flag.Parse()

	app, err := app.NewApp(configFilePath)
	if err != nil {
		log.Fatal(err)
	}

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
