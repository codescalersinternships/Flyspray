package app

import (
	"fmt"

	middleware "github.com/codescalersinternships/Flyspray/middlewares"
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

	return App{DB: database, router: gin.Default()}, nil
}

// App initializes the entire app
type App struct {
	DB     models.DBClient
	router *gin.Engine
}

// Run runs the server by seting the router and calling the internal setRoutes method
func (app *App) Run(port int) error {

	app.setRoutes()

	return app.router.Run(fmt.Sprintf(":%d", port))
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

	authUserGroup := authGroup.Group("/user")
	userGroup := app.router.Group("/user")

	{
		userGroup.POST("/signup", WrapFunc(app.signup))
		userGroup.POST("/signin", WrapFunc(app.signIn))
		userGroup.POST("/signup/verify", WrapFunc(app.verify))
		userGroup.POST("/refresh_token", WrapFunc(app.refreshToken))
	}

	{
		authUserGroup.PUT("", WrapFunc(app.updateUser))
		authUserGroup.GET("", WrapFunc(app.getUser))
	}

	component := app.router.Group("/component")
	{
		component.POST("", WrapFunc(app.createComponent))
		component.GET("/:id", WrapFunc(app.getComponent))
		component.DELETE("/:id", WrapFunc(app.deleteComponent))
		component.PUT("/:id", WrapFunc(app.updateComponent))
		component.GET("/filters", WrapFunc(app.getComponents))
	}

}
