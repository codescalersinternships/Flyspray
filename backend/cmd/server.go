package main

import (
	"flag"
	"log"

	"github.com/codescalersinternships/Flyspray/app"
)

func main() {
	// take port number and db filepath as flags
	var dbFilePath string
	var port int
	flag.StringVar(&dbFilePath, "d", "./flyspray.db", "Specify the filepath of sqlite database")
	flag.IntVar(&port, "p", 3000, "Specify the port number")
	flag.Parse()
	app, err := app.NewApp(dbFilePath)
	if err != nil {
		log.Fatal(err)
	}

	if err := app.Run(port); err != nil {
		log.Fatal(err)
	}
}
