package app

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/codescalersinternships/Flyspray/models"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

type createComponentInput struct {
	Name      string `json:"name" binding:"required"`
	ProjectID string `json:"project_id" binding:"required"`
}

type updateComponentInput struct {
	UserID string `json:"user_id" binding:"required"`
	Name   string `json:"name" binding:"required"`
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

func (a *App) updateComponent(ctx *gin.Context) (interface{}, Response) {
	id := ctx.Param("id")
	var input updateComponentInput

	if err := ctx.BindJSON(&input); err != nil {
		log.Error().Err(err).Send()
		return nil, BadRequest(errors.New("failed to read component data"))
	}

	userID, exists := ctx.Get("user_id")
	if !exists {
		return nil, UnAuthorized(errors.New("authentication is required"))
	}

	c, err := a.DB.GetComponent(id)

	if err == gorm.ErrRecordNotFound {
		log.Error().Err(err).Send()
		return nil, NotFound(errors.New("component is not found"))
	}
	if err != nil {
		log.Error().Err(err).Send()
		return nil, InternalServerError(errInternalServerError)
	}

	if userID != c.UserID {
		return nil, Forbidden(errors.New("have not access to update component"))
	}

	updatedComponent := models.Component{Name: input.Name}

	c, err = a.DB.GetComponentByName(input.Name)
	if err == nil && fmt.Sprint(c.ID) != id {
		return nil, BadRequest(errors.New("component name must be unique"))
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Error().Err(err).Send()
		return nil, InternalServerError(errInternalServerError)
	}

	err = a.DB.UpdateComponent(id, updatedComponent)

	if result.Error != nil {
		log.Println("error: Component not found")
		c.IndentedJSON(http.StatusNotFound,
			responseErr{Message: "Component not found"},
		)
		return
	}

	return ResponseMsg{
		Message: "component is updated successfully",
	}, Ok()
}

func (a *App) getComponent(ctx *gin.Context) (interface{}, Response) {
	id := ctx.Param("id")

	component, err := a.DB.GetComponent(id)

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

	return ResponseMsg{
		Message: "component is retrieved successfully",
		Data:    component,
	}, Ok()
}

func (a *App) deleteComponent(ctx *gin.Context) (interface{}, Response) {
	id := ctx.Param("id")

	userID, exists := ctx.Get("user_id")
	if !exists {
		return nil, UnAuthorized(errors.New("authentication is required"))
	}

	c, err := a.DB.GetComponent(id)

	if err == gorm.ErrRecordNotFound {
		log.Error().Err(err).Send()
		return nil, NotFound(errors.New("component is not found"))
	}
	if err != nil {
		log.Error().Err(err).Send()
		return nil, InternalServerError(errInternalServerError)
	}

	if fmt.Sprint(userID) != fmt.Sprint(c.UserID) {
		return nil, Forbidden(errors.New("have not access to delete component"))
	}

	err = a.DB.DeleteComponent(id)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Error().Err(err).Send()
		return nil, NotFound(errors.New("component is not found"))
	}
	if err != nil {
		log.Error().Err(err).Send()
		return nil, InternalServerError(errInternalServerError)
	}

	return ResponseMsg{
		Message: "component is deleted successfully",
	}, Ok()
}

func (a *App) getComponents(ctx *gin.Context) (interface{}, Response) {
	userId := ctx.Query("project_id")
	componentName := ctx.Query("name")
	creationDate := ctx.Query("after")

	components, err := a.DB.FilterComponents(userId, componentName, creationDate)

	if len(components) == 0 {
		log.Println("error: No components found for the specified project_id or name")
		c.IndentedJSON(http.StatusNotFound,
			responseErr{Message: "No components found for the specified project_id or name"},
		)
		return
	}

	return ResponseMsg{
		Message: "components are retrieved successfully",
		Data:    components,
	}, Ok()
}
