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
	a.setUserRoutes()

	return a.router.Run(fmt.Sprintf(":%d", port))
}

func (a *App) setUserRoutes() {
	a.router = gin.Default()

	userGroup := a.router.Group("/user")
	userGroup.POST("/signup", WrapFunc(a.Signup))
	userGroup.POST("/signin", WrapFunc(a.SignIn))
	userGroup.POST("/signup/verify", WrapFunc(a.Verify))
	userGroup.POST("/refresh_token", WrapFunc(a.RefreshToken))
	userGroup.PUT("", middleware.RequireAuth, WrapFunc(a.UpdateUser))
	userGroup.GET("", middleware.RequireAuth, WrapFunc(a.GetUser))
}
