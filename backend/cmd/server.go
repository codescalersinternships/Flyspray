package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/codescalersinternships/Flyspray/app"
	_ "github.com/codescalersinternships/Flyspray/docs"
)

var (
	version string
	commit  string
)

func main() {

	fmt.Println("Version:", version)
	fmt.Println("Commit:", commit)

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
