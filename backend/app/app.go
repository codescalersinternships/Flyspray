package app

import (
	"fmt"
	"net/http"

	"github.com/codescalersinternships/Flyspray/middlewares"
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

	a.router.POST("/signup", a.Signup)
	a.router.POST("/signin", a.SignIn)
	a.router.POST("/verify", a.Verify)
	a.router.GET("/testmiddleware", middlewares.RequireAuth, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, CustomResponse{
			Message: "HI",
		})
	})

	return a.router.Run(fmt.Sprintf(":%d", port))
}
