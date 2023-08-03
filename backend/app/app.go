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

// Run runs the server by seting the router and calling the internal registerHandlers method
func (app *App) Run(port int) error {

	app.router = gin.Default()

	app.registerHandlers()

	return app.router.Run(fmt.Sprintf(":%d", port))
}

func (app *App) registerHandlers() {

	comment := app.router.Group("/comment")
	{
		comment.POST("/", func(c *gin.Context) {
			app.CreateComment(c)

		})
		comment.GET("/:id", func(c *gin.Context) {
			app.GetComment(c)
		})
		comment.DELETE("/:id", func(c *gin.Context) {
			app.DeleteComment(c)
		})
		comment.GET("/filters", func(c *gin.Context) {
			app.ListComments(c)
		})
		comment.PUT("/:id", func(c *gin.Context) {
			app.UpdateComment(c)
		})

	}

}
