package app

import (
	"errors"
	"fmt"

	"github.com/codescalersinternships/Flyspray/models"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type createProjectInput struct {
	Name string `json:"name" binding:"required"`
}

type updateProjectInput struct {
	OwnerId uint   `json:"owner_id" binding:"required"`
	Name    string `json:"name" binding:"required"`
}

func (a *App) createProject(ctx *gin.Context) (interface{}, Response) {
	var input createProjectInput

	if err := ctx.BindJSON(&input); err != nil {
		log.Error().Err(err).Send()
		return nil, BadRequest(errors.New("failed to read project data"))
	}

	// TODO: get user id from authorization middleware and assign it to OwnerId
	newProject := models.Project{Name: input.Name, OwnerId: 10007} // 10007 is just a random number

	// check if project name exists before
	_, err := a.client.GetProjectByName(input.Name) // expected to return 'gorm.ErrRecordNotFound' if not exist
	if err == nil {                                 // project is found by name
		return nil, BadRequest(errors.New("project name must be unique"))
	}
	if err != gorm.ErrRecordNotFound { // there is some error other than not found
		log.Error().Err(err).Send()
		return nil, InternalServerError(errors.New("failed to create project"))
	}

	newProject, err = a.client.CreateProject(newProject)

	if err != nil {
		log.Error().Err(err).Send()
		return nil, InternalServerError(errors.New("failed to create project"))
	}

	return ResponseMsg{
		Message: "project is created successfully",
		Data:    newProject,
	}, Created()
}

func (a *App) updateProject(ctx *gin.Context) (interface{}, Response) {
	// TODO: get user id from authorization middleware and check if user has access to update the project
	id := ctx.Param("id")
	var input updateProjectInput

	if err := ctx.BindJSON(&input); err != nil {
		log.Error().Err(err).Send()
		return nil, BadRequest(errors.New("failed to read project data"))
	}

	updatedProject := models.Project{OwnerId: input.OwnerId, Name: input.Name}

	// check if project name exists before
	p, err := a.client.GetProjectByName(input.Name) // expected to return 'gorm.ErrRecordNotFound' if not exist
	if err == nil && fmt.Sprint(p.ID) != id {       // another project is found has the same name as the updated name
		return nil, BadRequest(errors.New("project name must be unique"))
	}
	if err != nil && err != gorm.ErrRecordNotFound { // there is some error and it is other than not found
		log.Error().Err(err).Send()
		return nil, InternalServerError(errors.New("failed to create project"))
	}

	err = a.client.UpdateProject(id, updatedProject)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Error().Err(err).Send()
		return nil, NotFound(errors.New("project is not found"))
	}
	if err != nil {
		log.Error().Err(err).Send()
		return nil, InternalServerError(errors.New("failed to update project"))
	}

	return ResponseMsg{
		Message: "project is updated successfully",
	}, Ok()
}

func (a *App) getProject(ctx *gin.Context) (interface{}, Response) {
	// TODO: add middleware to check if user is signed in
	id := ctx.Param("id")

	project, err := a.client.GetProject(id)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Error().Err(err).Send()
		return nil, NotFound(errors.New("project is not found"))
	}
	if err != nil {
		log.Error().Err(err).Send()
		return nil, InternalServerError(errors.New("failed to get project"))
	}

	return ResponseMsg{
		Message: "project is retrieved successfully",
		Data:    project,
	}, Ok()
}

func (a *App) getProjects(ctx *gin.Context) (interface{}, Response) {
	// TODO: add middleware to check if user is signed in
	userId := ctx.Query("userid")
	projectName := ctx.Query("name")
	creationDate := ctx.Query("after")

	projects, err := a.client.FilterProjects(userId, projectName, creationDate)

	if err != nil {
		log.Error().Err(err).Send()
		return nil, InternalServerError(errors.New("failed to get projects"))
	}

	return ResponseMsg{
		Message: "projects are retrieved successfully",
		Data:    projects,
	}, Ok()
}

func (a *App) deleteProject(ctx *gin.Context) (interface{}, Response) {
	// TODO: get user id from authorization middleware and check if user has access to delete the project
	id := ctx.Param("id")

	err := a.client.DeleteProject(id)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Error().Err(err).Send()
		return nil, NotFound(errors.New("project is not found"))
	}
	if err != nil {
		log.Error().Err(err).Send()
		return nil, InternalServerError(errors.New("failed to delete project"))
	}

	return ResponseMsg{
		Message: "project is deleted successfully",
	}, Ok()
}
