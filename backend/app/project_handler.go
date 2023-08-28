package app

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/codescalersinternships/Flyspray/models"
	"github.com/gin-gonic/gin"
	"github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

var errInternalServerError = errors.New("internal server error")

type createProjectInput struct {
	Name string `json:"name" binding:"required"`
}

type updateProjectInput struct {
	OwnerID string `json:"owner_id" binding:"required"`
	Name    string `json:"name" binding:"required"`
}

// createProject creates a new project
// @Summary Create a new project
// @Description Create a new project
// @Tags Projects
// @Accept json
// @Produce json
// @Param input body createProjectInput true "Project data"
// @Security ApiKeyAuth
// @Success 201 {object} ResponseMsg "Project is created successfully (Project details in the 'Data' field)"
// @Failure 400 {object} Response "Failed to read project data"
// @Failure 401 {object} Response "Authentication is required"
// @Failure 403 {object} Response "Access denied to create project"
// @Failure 409 {object} Response "Project name must be unique"
// @Failure 500 {object} Response "Internal server error"
// @Router /project [post]
func (a *App) createProject(ctx *gin.Context) (interface{}, Response) {
	var input createProjectInput

	if err := ctx.BindJSON(&input); err != nil {
		log.Error().Err(err).Send()
		return nil, BadRequest(errors.New("failed to read project data"))
	}

	userID, exists := ctx.Get("user_id")
	if !exists {
		return nil, UnAuthorized(errors.New("authentication is required"))
	}

	newProject := models.Project{Name: input.Name, OwnerID: fmt.Sprint(userID)}

	newProject, err := a.DB.CreateProject(newProject)

	var sqliteErr sqlite3.Error
	if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
		return nil, BadRequest(errors.New("project name must be unique"))
	}
	if err != nil {
		log.Error().Err(err).Send()
		return nil, InternalServerError(errInternalServerError)
	}
	// add project owner as a member in project
	member := models.Member{ProjectID: int(newProject.ID), Admin: true, UserID: newProject.OwnerID}
	err = a.DB.CreateNewMember(member)
	if err == gorm.ErrDuplicatedKey {
		log.Error().Err(err).Send()
		return nil, Forbidden(errors.New("member already exists"))
	}

	if err != nil {
		log.Error().Err(err).Send()
		return nil, BadRequest(errors.New("cannot create new member"))
	}
	return ResponseMsg{
		Message: "project is created successfully",
		Data:    newProject,
	}, Created()
}

// updateProject updates a project
// @Summary Update a project
// @Description Update an existing project
// @Tags Projects
// @Accept json
// @Produce json
// @Param id path string true "Project ID"
// @Param input body updateProjectInput true "Project data"
// @Security ApiKeyAuth
// @Success 200 {object} ResponseMsg "Project is updated successfully"
// @Failure 400 {object} Response "Failed to read project data"
// @Failure 401 {object} Response "Authentication is required"
// @Failure 403 {object} Response "Access denied to update project"
// @Failure 404 {object} Response "Project is not found"
// @Failure 409 {object} Response "Project name must be unique"
// @Failure 500 {object} Response "Internal server error"
// @Router /project/{id} [put]
func (a *App) updateProject(ctx *gin.Context) (interface{}, Response) {
	id := ctx.Param("id")
	var input updateProjectInput

	if err := ctx.BindJSON(&input); err != nil {
		log.Error().Err(err).Send()
		return nil, BadRequest(errors.New("failed to read project data"))
	}

	// check if he is project owner
	userID, exists := ctx.Get("user_id")
	if !exists {
		return nil, UnAuthorized(errors.New("authentication is required"))
	}

	p, err := a.DB.GetProject(id)

	if err == gorm.ErrRecordNotFound {
		log.Error().Err(err).Send()
		return nil, NotFound(errors.New("project is not found"))
	}
	if err != nil {
		log.Error().Err(err).Send()
		return nil, InternalServerError(errInternalServerError)
	}

	if fmt.Sprint(userID) != fmt.Sprint(p.OwnerID) {
		return nil, Forbidden(errors.New("have not access to update project"))
	}

	// proceed to update project
	convId, err := strconv.Atoi(id)
	if err != nil {
		log.Error().Err(err).Send()
		return nil, InternalServerError(errInternalServerError)
	}

	updatedProject := models.Project{ID: uint(convId), OwnerID: input.OwnerID, Name: input.Name}

	err = a.DB.UpdateProject(updatedProject)

	var sqliteErr sqlite3.Error
	if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
		return nil, BadRequest(errors.New("project name must be unique"))
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Error().Err(err).Send()
		return nil, NotFound(errors.New("project is not found"))
	}
	if err != nil {
		log.Error().Err(err).Send()
		return nil, InternalServerError(errInternalServerError)
	}

	return ResponseMsg{
		Message: "project is updated successfully",
	}, Ok()
}

// getProject retrieves a project
// @Summary Get a project
// @Description Get details of a project
// @Tags Projects
// @Produce json
// @Param id path string true "Project ID"
// @Success 200 {object} ResponseMsg "project is retrieved successfully (Project details in the 'Data' field)"
// @Failure 404 {object} Response "Project is not found"
// @Failure 500 {object} Response "Internal server error"
// @Router /project/{id} [get]
func (a *App) getProject(ctx *gin.Context) (interface{}, Response) {
	id := ctx.Param("id")

	project, err := a.DB.GetProject(id)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Error().Err(err).Send()
		return nil, NotFound(errors.New("project is not found"))
	}
	if err != nil {
		log.Error().Err(err).Send()
		return nil, InternalServerError(errInternalServerError)
	}

	return ResponseMsg{
		Message: "project is retrieved successfully",
		Data:    project,
	}, Ok()
}

// getProjects retrieves projects based on filters
// @Summary Get projects
// @Description Get a list of projects based on filters
// @Tags Projects
// @Produce json
// @Param userid query string false "User ID"
// @Param name query string false "Project name"
// @Param after query string false "Creation date (after)"
// @Success 200 {object} ResponseMsg "Projects are retrieved successfully (Projects details in the 'Data' field)"
// @Failure 500 {object} Response "Internal server error"
// @Router /project/filters [get]
func (a *App) getProjects(ctx *gin.Context) (interface{}, Response) {
	userId := ctx.Query("userid")
	projectName := ctx.Query("name")
	creationDate := ctx.Query("after")

	projects, err := a.DB.FilterProjects(userId, projectName, creationDate)

	if err != nil {
		log.Error().Err(err).Send()
		return nil, InternalServerError(errInternalServerError)
	}

	return ResponseMsg{
		Message: "projects are retrieved successfully",
		Data:    projects,
	}, Ok()
}

// deleteProject deletes a project
// @Summary Delete a project
// @Description Delete an existing project
// @Tags Projects
// @Produce json
// @Param id path string true "Project ID"
// @Security ApiKeyAuth
// @Success 200 {object} ResponseMsg "Project is deleted successfully"
// @Failure 401 {object} Response "Authentication is required"
// @Failure 403 {object} Response "Access denied to delete project"
// @Failure 404 {object} Response "Project is not found"
// @Failure 500 {object} Response "Internal server error"
// @Router /project/{id} [delete]
func (a *App) deleteProject(ctx *gin.Context) (interface{}, Response) {
	id := ctx.Param("id")

	// check if he is project owner
	userID, exists := ctx.Get("user_id")
	if !exists {
		return nil, UnAuthorized(errors.New("authentication is required"))
	}

	p, err := a.DB.GetProject(id)

	if err == gorm.ErrRecordNotFound {
		log.Error().Err(err).Send()
		return nil, NotFound(errors.New("project is not found"))
	}
	if err != nil {
		log.Error().Err(err).Send()
		return nil, InternalServerError(errInternalServerError)
	}

	if fmt.Sprint(userID) != fmt.Sprint(p.OwnerID) {
		return nil, Forbidden(errors.New("have not access to delete project"))
	}

	// proceed to delete project
	err = a.DB.DeleteProject(id)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Error().Err(err).Send()
		return nil, NotFound(errors.New("project is not found"))
	}
	if err != nil {
		log.Error().Err(err).Send()
		return nil, InternalServerError(errInternalServerError)
	}

	return ResponseMsg{
		Message: "project is deleted successfully",
	}, Ok()
}
