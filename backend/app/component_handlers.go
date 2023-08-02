package app

import (
	"log"
	"net/http"

	"github.com/codescalersinternships/Flyspray/models"

	"github.com/gin-gonic/gin"
)

// CreateComponent creates a new component
func (app *App) CreateComponent(c *gin.Context) {
	var component models.Component
	if err := c.ShouldBindJSON(&component); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if result := app.client.Client.Create(&component); result.Error != nil {
		log.Fatal(result.Error)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, component)
}

// GetComponentByID gets a component by its ID
func (app *App) GetComponentByID(c *gin.Context) {

	componentID := c.Param("id")
	component := models.Component{}
	result := app.client.Client.First(&component, componentID)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Component not found"})
		return
	}
	c.JSON(http.StatusOK, component)
}

// DeleteComponent deletes a component by its ID
func (app *App) DeleteComponent(c *gin.Context) {
	componentID := c.Param("id")

	component := models.Component{}

	db.Client.Delete(&component, componentID)

	c.JSON(http.StatusOK, gin.H{"message": "Component deleted successfully"})
}

// ListComponentsForProject gets all components for a project
func (app *App) ListComponentsForProject(c *gin.Context) {
	projectID := c.Query("project_id")
	name := c.Query("name")

	components := []models.Component{}
	query.Find(&components)

	c.JSON(http.StatusOK, components)
}
