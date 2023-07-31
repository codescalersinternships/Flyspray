package main

import (
	"flag"
	"log"

	"github.com/codescalersinternships/Flyspray/app"
	"github.com/codescalersinternships/Flyspray/models"
)

func main() {
	// take port number and db filepath as flags
	var dbFilePath string
	var port int

	flag.StringVar(&dbFilePath, "d", "./flyspray.db", "Specify the filepath of sqlite database")
	flag.IntVar(&port, "p", 8080, "Specify the port number")

	flag.Parse()

	client, err := models.NewDBClient(dbFilePath)
	if err != nil {
		log.Fatal(err)
	}

	app := app.NewApp(client)
	if err := app.Run(port); err != nil {
		log.Fatal(err)
	}
}
