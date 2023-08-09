package app

import (
	"fmt"
	"log"

	"github.com/codescalersinternships/Flyspray/models"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// NewApp is the factory of App
func NewApp(dbFilePath string) (App, error) {

	client, err := models.NewDBClient(dbFilePath)
	if err != nil {
		return App{}, err
	}

	if err := client.Migrate(); err != nil {
		log.Fatalf("error migrating tables %q", err)
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
	a.setRoutes()
	return a.router.Run(fmt.Sprintf(":%d", port))
}
func (a *App) setRoutes() {
	a.router = gin.Default()
	a.router.Use(cors.Default())
	memberRoutes := a.router.Group("/member")
	memberRoutes.POST("", WrapFunc(a.CreateNewMember))
	memberRoutes.GET("", WrapFunc(a.GetAllMembers))
	memberRoutes.PUT("/:id", WrapFunc(a.UpdateMemberOwnership))
}
