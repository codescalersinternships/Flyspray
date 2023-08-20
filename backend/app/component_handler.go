package app

import (
	"errors"
	"fmt"

	"github.com/codescalersinternships/Flyspray/models"
	"github.com/mattn/go-sqlite3"
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

func (a *App) createComponent(ctx *gin.Context) (interface{}, Response) {
	var input createComponentInput

	if err := ctx.BindJSON(&input); err != nil {
		log.Error().Err(err).Send()
		return nil, BadRequest(errors.New("failed to read component data"))
	}

	userID, exists := ctx.Get("user_id")
	if !exists {
		return nil, UnAuthorized(errors.New("authentication is required"))
	}

	newComponent := models.Component{Name: input.Name, ProjectID: input.ProjectID, UserID: userID.(string)}

	newComponent, err := a.DB.CreateComponent(newComponent)

	var sqliteErr sqlite3.Error
	if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
		return nil, BadRequest(errors.New("component name must be unique"))
	}
	if err != nil {
		log.Error().Err(err).Send()
		return nil, InternalServerError(errInternalServerError)
	}

	return ResponseMsg{
		Message: "component is created successfully",
		Data:    newComponent,
	}, Created()
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

	err = a.DB.UpdateComponent(id, updatedComponent)

	var sqliteErr sqlite3.Error
	if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
		return nil, BadRequest(errors.New("component name must be unique"))
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Error().Err(err).Send()
		return nil, NotFound(errors.New("component is not found"))
	}
	if err != nil {
		log.Error().Err(err).Send()
		return nil, InternalServerError(errInternalServerError)
	}

	return ResponseMsg{
		Message: "component is updated successfully",
	}, Ok()
}

func (a *App) getComponent(ctx *gin.Context) (interface{}, Response) {
	id := ctx.Param("id")

	component, err := a.DB.GetComponent(id)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Error().Err(err).Send()
		return nil, NotFound(errors.New("component is not found"))
	}
	if err != nil {
		log.Error().Err(err).Send()
		return nil, InternalServerError(errInternalServerError)
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

	if err != nil {
		log.Error().Err(err).Send()
		return nil, InternalServerError(errInternalServerError)
	}

	return ResponseMsg{
		Message: "components are retrieved successfully",
		Data:    components,
	}, Ok()
}
