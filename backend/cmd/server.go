package main

import (
	"flag"
	"log"

	"github.com/codescalersinternships/Flyspray/app"
	"github.com/codescalersinternships/Flyspray/internal"
)

func main() {

	err := internal.LoadEnv()
	if err != nil {
		log.Fatal(err)
	}
	// take port number and db filepath as flags
	var dbFilePath string
	var port int

	flag.StringVar(&dbFilePath, "d", "./flyspray.db", "Specify the filepath of sqlite database")
	flag.IntVar(&port, "p", 8080, "Specify the port number")

	flag.Parse()

	app, err := app.NewApp(dbFilePath)
	if err != nil {
		log.Fatal(err)
	}

	if err := app.Run(port); err != nil {
		log.Fatal(err)
	}

}
