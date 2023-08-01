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
	a.router.POST("project", func(ctx *gin.Context) {
		createProject(ctx, a.client)
	})
	a.router.GET("project/:id", func(ctx *gin.Context) {
		getProject(ctx, a.client)
	})
	a.router.GET("project/filters", func(ctx *gin.Context) {
		getProjects(ctx, a.client)
	})
	a.router.DELETE("project/:id", func(ctx *gin.Context) {
		deleteProject(ctx, a.client)
	})

	return a.router.Run(fmt.Sprintf(":%d", port))
}
