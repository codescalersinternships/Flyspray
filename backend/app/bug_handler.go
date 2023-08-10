package app

import (
	"errors"

	"github.com/codescalersinternships/Flyspray/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func (app *App) createBug(ctx *gin.Context) (interface{}, Response) {
	var bug models.Bug
	if err := ctx.BindJSON(&bug); err != nil {
		return nil, BadRequest(errors.New("failed to read data"))
	}

	// Create a new instance of the validator
	validate := validator.New()
	// Validate the bug struct
	if err := validate.Struct(bug); err != nil {
		return nil, BadRequest(errors.New("validation error: " + err.Error()))
	}

	// TODO: add middleware to check if user is signed in

	newBug, err := app.client.CreateNewBug(bug)
	if err != nil {
		return nil, InternalServerError(errors.New("failed to create bug"))
	}

	return ResponseMsg{
		Message: "bug is created successfully",
		Data:    newBug,
	}, Created()
}

func (a *App) getBugs(ctx *gin.Context) (interface{}, Response) {
	// TODO: add middleware to check if user is signed in

	// filters
	var (
		userId      = ctx.Query("user_id")
		category    = ctx.Query("category")
		status      = ctx.Query("status")
		componentId = ctx.Query("component_id")
	)

	bugs, err := a.client.FilterBugs(userId, category, status, componentId)

	if err != nil {
		return nil, InternalServerError(errors.New("failed to get bugs"))
	}

	return ResponseMsg{
		Message: "bugs is retrieved successfully",
		Data:    bugs,
	}, Ok()
}

func (app *App) getSpecificBug(ctx *gin.Context) (interface{}, Response) {
	// TODO: add middleware to check if user is signed in
	id := ctx.Param("id")

	bug, err := app.client.GetSpecificBug(id)

	if err == gorm.ErrRecordNotFound {
		return nil, NotFound(errors.New("bug is not found"))
	}

	if err != nil {
		return nil, InternalServerError(errors.New("failed to get bug"))
	}

	return ResponseMsg{
		Message: "bug is retrieved successfully",
		Data:    bug,
	}, Ok()
}

func (app *App) updateBug(ctx *gin.Context) (interface{}, Response) {
	// TODO: add middleware to check if user is signed in
	id := ctx.Param("id")

	bug := models.Bug{}

	if result := app.client.Client.First(&bug, id); result.Error != nil {
		return nil, NotFound(errors.New("bug is not found"))
	}

	updatedBug := bug
	// Create a new instance of the validator
	validate := validator.New()
	// Validate the bug struct
	if err := validate.Struct(bug); err != nil {
		return nil, BadRequest(errors.New("validation error: " + err.Error()))
	}

	if err := ctx.BindJSON(&updatedBug); err != nil {
		return nil, BadRequest(errors.New("failed to read data"))
	}

	result := app.client.UpdateBug(id, updatedBug)
	if result.Error != nil {
		return nil, InternalServerError(errors.New("field to update bug"))
	}

	return ResponseMsg{
		Message: "bug is updated successfully",
		Data:    updatedBug,
	}, Ok()
}

func (app *App) deleteBug(ctx *gin.Context) (interface{}, Response) {
	// TODO: add middleware to check if user is signed in
	id := ctx.Param("id")

	if id == "" {
		return nil, BadRequest(errors.New("bug ID is required"))
	}

	if err := app.client.DeleteBug(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, NotFound(errors.New("bug is not found"))
		} else {
			return nil, InternalServerError(errors.New("failed to delete bug"))
		}
	}

	return ResponseMsg{
		Message: "bug is deleted successfully",
	}, Ok()
}
