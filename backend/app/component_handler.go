package app

import (
	"errors"
	"strconv"

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
	Name string `json:"name" binding:"required"`
}

// createComponent creates a new component
// @Summary Create a component
// @Description Create a new component for a project
// @Tags Components
// @Accept json
// @Produce json
// @Param input body createComponentInput true "Component data"
// @Security ApiKeyAuth
// @Success 201 {object} ResponseMsg "component is created successfully (Component details in the 'Data' field)"
// @Failure 400 {object} Response "Failed to read component data"
// @Failure 401 {object} Response "Authentication is required"
// @Failure 403 {object} Response "Do not have access to create component"
// @Failure 500 {object} Response "Internal server error"
// @Router /component [post]
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

	projectIdInt, err := strconv.Atoi(input.ProjectID)
	if err != nil {
		log.Error().Err(err).Send()
		return nil, BadRequest(errors.New("invalid project id"))
	}

	err = a.DB.CheckMembers(projectIdInt, userID.(string))

	if err != nil {
		return nil, Forbidden(errors.New("have not access to create component"))
	}

	newComponent := models.Component{Name: input.Name, ProjectID: input.ProjectID, UserID: userID.(string)}

	newComponent, err = a.DB.CreateComponent(newComponent)

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

// updateComponent updates an existing component
// @Summary Update a component
// @Description Update an existing component
// @Tags Components
// @Accept json
// @Produce json
// @Param id path string true "Component ID"
// @Param input body updateComponentInput true "Component data"
// @Security ApiKeyAuth
// @Success 200 {object} ResponseMsg "component is updated successfully"
// @Failure 400 {object} Response "Failed to read component data"
// @Failure 401 {object} Response "Authentication is required"
// @Failure 403 {object} Response "Do not have access to update component"
// @Failure 404 {object} Response "Component is not found"
// @Failure 500 {object} Response "Internal server error"
// @Router /component/{id} [put]
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

	project, err := a.DB.GetProject(c.ProjectID)
	if err != nil {
		log.Error().Err(err).Send()
		return nil, InternalServerError(errInternalServerError)
	}

	if userID != c.UserID && userID != project.OwnerID {
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

// @Summary Get a component
// @Description Get a component by ID
// @Tags Components
// @Produce json
// @Param id path string true "Component ID"
// @Success 200 {object} ResponseMsg "component is retrieved successfully (Component details in the 'Data' field)"
// @Failure 404 {object} Response "Component is not found"
// @Failure 500 {object} Response "Internal server error"
// @Router /component/{id} [get]
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

// deleteComponent deletes a component by ID
// @Summary Delete a component
// @Description Delete a component by ID
// @Tags Components
// @Produce json
// @Param id path string true "Component ID"
// @Security ApiKeyAuth
// @Success 200 {object} ResponseMsg "Component is deleted successfully"
// @Failure 401 {object} Response "Authentication is required"
// @Failure 403 {object} Response "Do not have access to delete component"
// @Failure 404 {object} Response "Component is not found"
// @Failure 500 {object} Response "Internal server error"
// @Router /component/{id} [delete]
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

	project, err := a.DB.GetProject(c.ProjectID)
	if err != nil {
		log.Error().Err(err).Send()
		return nil, InternalServerError(errInternalServerError)
	}

	if userID != c.UserID && userID != project.OwnerID {
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

// getComponents retrieves a list of components based on filters
// @Summary Get components
// @Description Get a list of components based on filters
// @Tags Components
// @Produce json
// @Param project_id query string false "Project ID"
// @Param name query string false "Component name"
// @Param after query string false "Creation date (after)"
// @Success 200 {object} ResponseMsg "Components are retrieved successfully (Components details in the 'Data' field)"
// @Failure 500 {object} Response "Internal server error"
// @Router /component/filters [get]
func (a *App) getComponents(ctx *gin.Context) (interface{}, Response) {
	projectId := ctx.Query("project_id")
	componentName := ctx.Query("name")
	creationDate := ctx.Query("after")

	components, err := a.DB.FilterComponents(projectId, componentName, creationDate)

	if err != nil {
		log.Error().Err(err).Send()
		return nil, InternalServerError(errInternalServerError)
	}

	return ResponseMsg{
		Message: "components are retrieved successfully",
		Data:    components,
	}, Ok()
}
