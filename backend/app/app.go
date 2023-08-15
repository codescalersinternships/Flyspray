package app

import (
	"fmt"

	"github.com/codescalersinternships/Flyspray/models"
	"github.com/gin-gonic/gin"
)

// NewApp is the factory of App
func NewApp(dbFilePath string) (App, error) {

	database, err := models.NewDBClient(dbFilePath)
	if err != nil {
		return App{}, err
	}

	if err := database.Migrate(); err != nil {
		return App{}, err
	}

	return App{db: database}, nil
}

// App initializes the entire app
type App struct {
	db     models.DBClient
	router *gin.Engine
}

// Run runs the server by seting the router and calling the internal registerHandlers method
func (app *App) Run(port int) error {

	app.router = gin.Default()

	app.setRoutes()

	return app.router.Run(fmt.Sprintf(":%d", port))
}

func (app *App) setRoutes() {

	comment := app.router.Group("/comment")
	{
		comment.POST("/", WrapFunc(app.createComment))

		comment.GET("/:id", WrapFunc(app.getComment))

		comment.DELETE("/:id", WrapFunc(app.deleteComment))

		comment.GET("/filters", WrapFunc(app.listComments))

		comment.PUT("/:id", WrapFunc(app.updateComment))
	}

}
