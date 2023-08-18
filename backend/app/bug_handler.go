package app

import (
	"errors"

	"github.com/codescalersinternships/Flyspray/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

var errInternalServerError = errors.New("internal server error")

type createBugInput struct {
	Summary     string `json:"summary"`
	ComponentID int    `json:"component_id" validate:"required"`
}

type updateBugInput struct {
	Summary string `json:"summary"`
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

	// TODO: add middleware to check if user is signed in

	newBug, err := app.client.CreateNewBug(models.Bug{UserID: userID.(string), ComponentID: bug.ComponentID})
	if err != nil {
		log.Error().Err(err).Send()
		return nil, InternalServerError(errInternalServerError)
	}

	return ResponseMsg{
		Message: "bug is created successfully",
		Data:    newBug,
	}, Created()
}

func (a *App) getbug(ctx *gin.Context) (interface{}, Response) {
	// TODO: add middleware to check if user is signed in

	// filters
	var (
		category    = ctx.Query("category")
		status      = ctx.Query("status")
		componentId = ctx.Query("component_id")
	)

	bug, err := a.client.Filterbug(category, status, componentId)
	if err != nil {
		log.Error().Err(err).Send()
		return nil, InternalServerError(errInternalServerError)
	}

	return ResponseMsg{
		Message: "bug is retrieved successfully",
		Data:    bug,
	}, Ok()
}

func (app *App) getSpecificBug(ctx *gin.Context) (interface{}, Response) {
	// TODO: add middleware to check if user is signed in
	id := ctx.Param("id")

	if id == "" {
		return nil, BadRequest(errors.New("bug id is required"))
	}

	bug, err := app.client.GetSpecificBug(id)

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

	_, exists := ctx.Get("user_id")
	if !exists {
		return nil, UnAuthorized(errors.New("authentication is required"))
	}

	_, err := app.client.GetSpecificBug(id)

	if err == gorm.ErrRecordNotFound {
		log.Error().Err(err).Send()
		return nil, NotFound(errors.New("bug is not found"))
	}

	if err != nil {
		log.Error().Err(err).Send()
		return nil, InternalServerError(errInternalServerError)
	}

	// proceed to update bug
	updatedBug := models.Bug{Summary: input.Summary}

	err = app.client.UpdateBug(id, updatedBug)
	if err != nil {
		log.Error().Err(err).Send()
		return nil, InternalServerError(errors.New("field to update bug"))
	}

	return ResponseMsg{
		Message: "bug is updated successfully",
	}, Ok()
}

func (app *App) deleteBug(ctx *gin.Context) (interface{}, Response) {
	// TODO: add middleware to check if user is signed in
	id := ctx.Param("id")

	if id == "" {
		return nil, BadRequest(errors.New("bug ID is required"))
	}

	userId, exists := ctx.Get("user_id")
	if !exists {
		return nil, UnAuthorized(errors.New("authentication is required"))
	}

	bug, err := app.client.GetSpecificBug(id)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Error().Err(err).Send()
		return nil, NotFound(errors.New("bug is not found"))
	}

	if err != nil {
		log.Error().Err(err).Send()
		return nil, InternalServerError(errInternalServerError)
	}

	if userId != bug.UserID {
		return nil, Forbidden(errors.New("you have no access to delete the bug"))
	}

	err = app.client.DeleteBug(id)

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
