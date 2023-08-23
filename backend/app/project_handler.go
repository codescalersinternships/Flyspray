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
