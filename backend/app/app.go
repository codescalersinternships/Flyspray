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

	a.setRoutes()

	return a.router.Run(fmt.Sprintf(":%d", port))
}

func (a *App) setRoutes() {
	project := a.router.Group("/project")
	project.POST("", a.createProject)
	project.GET("/filters", a.getProjects)
	project.GET("/:id", a.getProject)
	project.PUT("/:id", a.updateProject)
	project.DELETE("/:id", a.deleteProject)
}
