package app

import (
	"fmt"

	"github.com/codescalersinternships/Flyspray/models"
	"github.com/gin-gonic/gin"
)

// NewApp is the factory of App
func NewApp(dbFilePath string) (App, error) {

	client, err := models.NewDBClient(dbFilePath)
	if err != nil {
		return App{}, err
	}

	if err := client.Migrate(); err != nil {
		return App{}, err
	}

	return App{client: client}, nil
}

// App initializes the entire app
type App struct {
	client models.DBClient
	router *gin.Engine
}

// Run runs server
func (a *App) Run(port int) error {
	a.router = gin.Default()

	// set routes here

	return a.router.Run(fmt.Sprintf(":%d", port))
}
