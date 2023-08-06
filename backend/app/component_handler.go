package app

import (
	"log"
	"net/http"
	"time"

	"github.com/codescalersinternships/Flyspray/models"

	"github.com/gin-gonic/gin"
)

type responseErr struct {
	Message string `json:"message"`
}

type responseOk struct {
	Message string             `json:"message"`
	Data    []models.Component `json:"data"`
}

// CreateComponent creates a new component
func (app *App) CreateComponent(c *gin.Context) {
	var component models.Component
	if err := c.ShouldBindJSON(&component); err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusBadRequest,
			responseErr{Message: err.Error()},
		)
		return
	}
	component.CreatedAt = time.Now()

	if result := app.client.Client.Create(&component); result.Error != nil {
		log.Println("error: Failed to Create the component")
		c.IndentedJSON(http.StatusInternalServerError,
			responseErr{Message: "Failed to Create the component"},
		)
		return
	}
	c.IndentedJSON(http.StatusCreated, responseOk{
		Message: "project created successfully",
		Data:    []models.Component{component},
	})
}

// GetComponentByID gets a component by its ID
func (app *App) GetComponentByID(c *gin.Context) {

	componentID := c.Param("id")
	component := models.Component{}
	result := app.client.Client.First(&component, componentID)

	if result.Error != nil {
		log.Println("error: Component not found")
		c.IndentedJSON(http.StatusNotFound,
			responseErr{Message: "Component not found"},
		)
		return
	}
	c.IndentedJSON(http.StatusOK, responseOk{
		Message: "project retrieved successfully",
		Data:    []models.Component{component},
	})
}

// DeleteComponent deletes a component by its ID
func (app *App) DeleteComponent(c *gin.Context) {
	componentID := c.Param("id")

	component := models.Component{}

	if result := app.client.Client.First(&component, componentID); result.Error != nil {
		log.Println("error: Component not found")
		c.IndentedJSON(http.StatusNotFound,
			responseErr{Message: "Component not found"},
		)
		return
	}
	if result := app.client.Client.Delete(&component, componentID); result.Error != nil {
		log.Println("error: Failed to delete the component")
		c.IndentedJSON(http.StatusInternalServerError,
			responseErr{Message: "Failed to delete the component"},
		)
		return
	}

	c.IndentedJSON(http.StatusOK, responseOk{
		Message: "project deleted successfully",
		Data:    []models.Component{component},
	})
}

// UpdateComponent updates an existing component
func (app *App) UpdateComponent(c *gin.Context) {

	componentID := c.Param("id")
	updatedComponent := models.Component{}

	component := models.Component{}

	if result := app.client.Client.First(&component, componentID); result.Error != nil {
		log.Println("error: Component not found")
		c.IndentedJSON(http.StatusNotFound,
			responseErr{Message: "Component not found"},
		)
		return
	}

	if err := c.ShouldBindJSON(&updatedComponent); err != nil {
		log.Println(err)
		c.IndentedJSON(http.StatusBadRequest,
			responseErr{Message: err.Error()},
		)
		return
	}

	result := app.client.Client.Model(&models.Component{}).Where("id = ?", componentID).Updates(updatedComponent)
	if result.Error != nil {
		log.Println(result.Error.Error())
		c.IndentedJSON(http.StatusInternalServerError,
			responseErr{Message: result.Error.Error()},
		)
		return
	}

	c.IndentedJSON(http.StatusOK, responseOk{
		Message: "project updated successfully",
		Data:    []models.Component{component},
	})
}

// ListComponentsForProject gets all components for a project
func (app *App) ListComponentsForProject(c *gin.Context) {
	projectID := c.Query("project_id")
	name := c.Query("name")

	if projectID == "" {
		fmt.Println(projectID)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing or invalid project_id"})
		return
	}

	query := app.client.Client.Where("project_id = ?", projectID)

	if name != "" {
		query = query.Where("name = ?", name)
	}

	components := []models.Component{}
	query.Find(&components)

	if len(components) == 0 {
		log.Println("error: No components found for the specified project_id or name")
		c.IndentedJSON(http.StatusNotFound,
			responseErr{Message: "No components found for the specified project_id or name"},
		)
		return
	}

	c.JSON(http.StatusOK, components)
}
