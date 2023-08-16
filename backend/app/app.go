package app

import (
	"fmt"

	middleware "github.com/codescalersinternships/Flyspray/middlewares"
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

	// set routes here
	a.setRoutes()

	return a.router.Run(fmt.Sprintf(":%d", port))
}

func (a *App) setRoutes() {
	a.router = gin.Default()

	authUserGroup := a.router.Group("/user")
	userGroup := authUserGroup.Group("")
	
	
	userGroup.POST("/signup", WrapFunc(a.signup))
	userGroup.POST("/signin", WrapFunc(a.signIn))
	userGroup.POST("/signup/verify", WrapFunc(a.verify))
	userGroup.POST("/refresh_token", WrapFunc(a.refreshToken))
	
	authUserGroup.Use(middleware.RequireAuth)
	
	authUserGroup.PUT("", WrapFunc(a.updateUser))
	authUserGroup.GET("", WrapFunc(a.getUser))
}
