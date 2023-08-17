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

	return App{DB: database}, nil
}

// App initializes the entire app
type App struct {
	DB     models.DBClient
	router *gin.Engine
}

// Run runs the server by seting the router and calling the internal setRoutes method
func (app *App) Run(port int) error {

	app.router = gin.Default()

	app.setRoutes()

	return app.router.Run(fmt.Sprintf(":%d", port))
}

func (app *App) setRoutes() {
	project := app.router.Group("/project")
	{
		project.POST("", WrapFunc(app.createProject))
		project.GET("/filters", WrapFunc(app.getProjects))
		project.GET("/:id", WrapFunc(app.getProject))
		project.PUT("/:id", WrapFunc(app.updateProject))
		project.DELETE("/:id", WrapFunc(app.deleteProject))

	}

	comment := app.router.Group("/comment")
	{
		comment.POST("", WrapFunc(app.createComment))
		comment.GET("/:id", WrapFunc(app.getComment))
		comment.DELETE("/:id", WrapFunc(app.deleteComment))
		comment.GET("/filters", WrapFunc(app.listComments))
		comment.PUT("/:id", WrapFunc(app.updateComment))
	}
	memberRoutes := app.router.Group("/member")
	{
		memberRoutes.POST("", WrapFunc(app.createNewMember))
		memberRoutes.PUT("/:id", WrapFunc(app.updateMemberOwnership))
		memberRoutes.GET("/:project_id", WrapFunc(app.getMembersInProject))
	}

	component := app.router.Group("/component")
	{
		component.POST("/", WrapFunc(app.createComponent))
		component.GET("/:id", WrapFunc(app.getComponent))
		component.DELETE("/:id", WrapFunc(app.deleteComment))
		component.PUT("/:id", WrapFunc(app.updateComment))
		component.GET("/filters", WrapFunc(app.getComponents))
	}

	component := app.router.Group("/component")
	{
		component.POST("/", WrapFunc(app.createComponent))
		component.GET("/:id", WrapFunc(app.getComponent))
		component.DELETE("/:id", WrapFunc(app.deleteComment))
		component.PUT("/:id", WrapFunc(app.updateComment))
		component.GET("/filters", WrapFunc(app.getComponents))
	}

}
