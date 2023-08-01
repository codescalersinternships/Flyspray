package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/codescalersinternships/Flyspray/models"

	"github.com/gin-gonic/gin"
)

// CreateComponent creates a new component
func CreateComponent(c *gin.Context, db models.DBClient) {
	var component models.Component
	if err := c.ShouldBindJSON(&component); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if result := db.Client.Create(&component); result.Error != nil {
		log.Fatal(result.Error)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, component)
}

// GetComponentByID gets a component by its ID
func GetComponentByID(c *gin.Context, db models.DBClient) {

	componentID := c.Param("id")
	component := models.Component{}
	result := db.Client.First(&component, componentID)

	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Component not found"})
		return
	}
	c.JSON(http.StatusOK, component)
}

// DeleteComponent deletes a component by its ID
func DeleteComponent(c *gin.Context, db models.DBClient) {
	componentID := c.Param("id")

	component := models.Component{}

	if result := db.Client.First(&component, componentID); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Component not found"})
		return
	}
	if result := db.Client.Delete(&component, componentID); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete the component"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Component deleted successfully"})
}

// ListComponentsForProject gets all components for a project
func ListComponentsForProject(c *gin.Context, db models.DBClient) {
	projectID := c.Query("project_id")

	fmt.Println(projectID)
	if projectID == "" {
		fmt.Println(projectID)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing or invalid project_id"})
		return
	}

	components := []models.Component{}
	db.Client.Where("project_id = ?", projectID).Find(&components)

	if len(components) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No components found for the specified project_id"})
		return
	}

	c.JSON(http.StatusOK, components)
}
