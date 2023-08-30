package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/codescalersinternships/Flyspray/app"
	_ "github.com/codescalersinternships/Flyspray/docs"
)

var (
	version string
	commit  string
)

func main() {

	// if the user enters version in the run command then the latest version & commit of the app will be printed
	if len(os.Args) > 1 && os.Args[1] == "version" {
		// Print the version and commit information
		fmt.Println("Version:", version)
		fmt.Println("Commit:", commit)
		return
	}

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
