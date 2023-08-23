package app

import (
	"errors"
	"fmt"
	"time"

	"github.com/codescalersinternships/Flyspray/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type createBugInput struct {
	Summary     string `json:"summary"`
	ComponentID int    `json:"component_id" validate:"required"`
	Category    string `json:"category"`
	Severity    string `json:"severity"`
	Status      string `json:"status"`
}

type updateBugInput struct {
	Summary  string `json:"summary"`
	Category string `json:"category"`
	Severity string `json:"severity"`
	Status   string `json:"status"`
}

func (app *App) createBug(ctx *gin.Context) (interface{}, Response) {
	var bug createBugInput

	if err := ctx.BindJSON(&bug); err != nil {
		log.Error().Err(err).Send()
		return nil, BadRequest(errors.New("failed to read data"))
	}

	userID, exists := ctx.Get("user_id")
	if !exists {
		return nil, UnAuthorized(errors.New("authentication is required"))
	}

	// Create a new instance of the validator
	validate := validator.New()
	// Validate the bug struct
	if err := validate.Struct(bug); err != nil {
		log.Error().Err(err).Send()
		return nil, BadRequest(errors.New("validation error: " + err.Error()))
	}

	newBug, err := app.DB.CreateNewBug(models.Bug{UserID: userID.(string), ComponentID: bug.ComponentID, Summary: bug.Summary, Category: bug.Category, Status: bug.Status, Severity: bug.Severity, Opened: true, OpenedAt: time.Now()})
	if err != nil {
		log.Error().Err(err).Send()
		return nil, InternalServerError(errInternalServerError)
	}

	return ResponseMsg{
		Message: "bug is created successfully",
		Data:    newBug,
	}, Created()
}

func (a *App) getbugs(ctx *gin.Context) (interface{}, Response) {

	_, exists := ctx.Get("user_id")
	if !exists {
		return nil, UnAuthorized(errors.New("authentication is required"))
	}

	// filters
	var (
		category    = ctx.Query("category")
		status      = ctx.Query("status")
		componentId = ctx.Query("component_id")
	)

	bug, err := a.DB.Filterbug(category, status, componentId)
	if err != nil {
		log.Error().Err(err).Send()
		return nil, InternalServerError(errInternalServerError)
	}

	return ResponseMsg{
		Message: "bug is retrieved successfully",
		Data:    bug,
	}, Ok()
}

func (app *App) getBug(ctx *gin.Context) (interface{}, Response) {

	_, exists := ctx.Get("user_id")
	if !exists {
		return nil, UnAuthorized(errors.New("authentication is required"))
	}

	id := ctx.Param("id")

	if id == "" {
		return nil, BadRequest(errors.New("bug id is required"))
	}

	bug, err := app.DB.GetBug(id)

	if err == gorm.ErrRecordNotFound {
		log.Error().Err(err).Send()
		return nil, NotFound(errors.New("bug is not found"))
	}

	if err != nil {
		log.Error().Err(err).Send()
		return nil, InternalServerError(errInternalServerError)
	}

	return ResponseMsg{
		Message: "bug is retrieved successfully",
		Data:    bug,
	}, Ok()
}

func (app *App) updateBug(ctx *gin.Context) (interface{}, Response) {
	var input updateBugInput

	if err := ctx.BindJSON(&input); err != nil {
		log.Error().Err(err).Send()
		return nil, BadRequest(errors.New("failed to read bug data"))
	}

	id := ctx.Param("id")
	if id == "" {
		return nil, BadRequest(errors.New("bug id is required"))
	}

	userID, exists := ctx.Get("user_id")
	if !exists {
		return nil, UnAuthorized(errors.New("authentication is required"))
	}

	bug, err := app.DB.GetBug(id)

	if err == gorm.ErrRecordNotFound {
		log.Error().Err(err).Send()
		return nil, NotFound(errors.New("bug is not found"))
	}

	if err != nil {
		log.Error().Err(err).Send()
		return nil, InternalServerError(errInternalServerError)
	}

	// check if user is autherized to update the bug or not
	if fmt.Sprint(userID) != fmt.Sprint(bug.UserID) {
		return nil, Forbidden(errors.New("you have no access to update the bug"))
	}

	// proceed to update bug
	updatedBug := models.Bug{Summary: input.Summary, Category: input.Category, Severity: input.Severity, Status: input.Status, UpdatedAt: time.Now()}

	err = app.DB.UpdateBug(id, updatedBug)
	if err != nil {
		log.Error().Err(err).Send()
		return nil, InternalServerError(errors.New("field to update bug"))
	}

	return ResponseMsg{
		Message: "bug is updated successfully",
	}, Ok()
}

func (app *App) deleteBug(ctx *gin.Context) (interface{}, Response) {

	id := ctx.Param("id")

	if id == "" {
		return nil, BadRequest(errors.New("bug ID is required"))
	}

	userId, exists := ctx.Get("user_id")
	if !exists {
		return nil, UnAuthorized(errors.New("authentication is required"))
	}

	bug, err := app.DB.GetBug(id)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Error().Err(err).Send()
		return nil, NotFound(errors.New("bug is not found"))
	}

	if err != nil {
		log.Error().Err(err).Send()
		return nil, InternalServerError(errInternalServerError)
	}

	component, err := app.DB.GetComponent(id)

	if err == gorm.ErrRecordNotFound {
		log.Error().Err(err).Send()
		return nil, NotFound(errors.New("component is not found"))
	}
	if err != nil {
		log.Error().Err(err).Send()
		return nil, InternalServerError(errInternalServerError)
	}

	project, err := app.DB.GetProject(component.ProjectID)
	if err != nil {
		log.Error().Err(err).Send()
		return nil, InternalServerError(errInternalServerError)
	}

	members, err := app.DB.GetMembersInProject(int(project.ID))
	if err != nil {
		log.Error().Err(err).Send()
		return nil, InternalServerError(errInternalServerError)
	}

	adminMember := false

	// loop over the members to check admin state
	for _, member := range members {

		if userId == member.ID && member.Admin {

			adminMember = true
			break

		}
	}

	if userId != bug.UserID && userId != project.OwnerID && !adminMember {
		return nil, Forbidden(errors.New("you have no access to delete the bug"))
	}

	err = app.DB.DeleteBug(id)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Error().Err(err).Send()
		return nil, NotFound(errors.New("bug is not found"))

	}

	if err != nil {
		log.Error().Err(err).Send()
		return nil, InternalServerError(errInternalServerError)
	}

	return ResponseMsg{
		Message: "bug is deleted successfully",
	}, Ok()
}
