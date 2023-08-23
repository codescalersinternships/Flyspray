package app

import (
	"fmt"

	"github.com/codescalersinternships/Flyspray/internal"
	middleware "github.com/codescalersinternships/Flyspray/middlewares"
	"github.com/codescalersinternships/Flyspray/models"
	"github.com/gin-gonic/gin"
)

// NewApp is the factory of App
func NewApp(configFilePath string) (App, error) {
	config, err := internal.ReadConfigFile(configFilePath)
	if err != nil {
		return App{}, err
	}

	database, err := models.NewDBClient(config.DB.File)
	if err != nil {
		return App{}, err
	}

	if err := database.Migrate(); err != nil {
		return App{}, err
	}

	return App{config: config, DB: database, router: gin.Default()}, nil
}

// App initializes the entire app
type App struct {
	config internal.Configuration
	DB     models.DBClient
	router *gin.Engine
}

// Run runs the server by setting the router and calling the internal setRoutes method
func (app *App) Run() error {
	app.setRoutes()

	return app.router.Run(fmt.Sprintf(":%d", app.config.Server.Port))
}

func (app *App) setRoutes() {
	authGroup := app.router.Group("")
	authGroup.Use(middleware.RequireAuth(""))

	project := authGroup.Group("/project")
	{
		project.POST("", WrapFunc(app.createProject))
		project.GET("/filters", WrapFunc(app.getProjects))
		project.GET("/:id", WrapFunc(app.getProject))
		project.PUT("/:id", WrapFunc(app.updateProject))
		project.DELETE("/:id", WrapFunc(app.deleteProject))

	}

	comment := authGroup.Group("/comment")
	{
		comment.POST("", WrapFunc(app.createComment))
		comment.GET("/:id", WrapFunc(app.getComment))
		comment.DELETE("/:id", WrapFunc(app.deleteComment))
		comment.GET("/filters", WrapFunc(app.listComments))
		comment.PUT("/:id", WrapFunc(app.updateComment))
	}

	memberRoutes := authGroup.Group("/member")
	{
		memberRoutes.POST("", WrapFunc(app.createNewMember))
		memberRoutes.PUT("/:id", WrapFunc(app.updateMemberOwnership))
		memberRoutes.GET("/:project_id", WrapFunc(app.getMembersInProject))
	}

	userGroup := app.router.Group("/user")
	{
		userGroup.POST("/signup", WrapFunc(app.signup))
		userGroup.POST("/signin", WrapFunc(app.signIn))
		userGroup.POST("/signup/verify", WrapFunc(app.verify))
		userGroup.POST("/refresh_token", WrapFunc(app.refreshToken))
	}

	authUserGroup := authGroup.Group("/user")
	{
		authUserGroup.PUT("", WrapFunc(app.updateUser))
		authUserGroup.GET("", WrapFunc(app.getUser))
	}

	component := authGroup.Group("/component")
	{
		component.POST("", WrapFunc(app.createComponent))
		component.GET("/:id", WrapFunc(app.getComponent))
		component.DELETE("/:id", WrapFunc(app.deleteComponent))
		component.PUT("/:id", WrapFunc(app.updateComponent))
		component.GET("/filters", WrapFunc(app.getComponents))
	}
}
