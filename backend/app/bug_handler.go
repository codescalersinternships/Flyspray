package app

import (
	"errors"
	"strconv"
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

// createBug creates a new bug
// @Summary Create a bug
// @Description Creates a new bug with the provided input
// @Tags bugs
// @Accept json
// @Produce json
// @Param bug body createBugInput true "Bug data"
// @Security ApiKeyAuth
// @Success 201 {object} ResponseMsg "bug is created successfully (Bug details in the 'Data' field)"
// @Failure 400 {object} Response "Bad Request"
// @Failure 401 {object} Response "Unauthorized"
// @Failure 500 {object} Response "Internal Server Error"
// @Router /bug [post]
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

// getbugs retrieves bugs based on the provided filters
// @Summary Get bugs
// @Description Retrieves bugs based on the provided filters
// @Tags bugs
// @Produce json
// @Param category query string false "Bug category"
// @Param status query string false "Bug status"
// @Param component_id query string false "Component ID"
// @Security ApiKeyAuth
// @Success 200 {object} ResponseMsg "bugs are retrieved successfully (Bugs details in the 'Data' field)"
// @Failure 401 {object} Response "Unauthorized"
// @Failure 500 {object} Response "Internal Server Error"
// @Router /bug/filters [get]
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

	bugs, err := a.DB.Filterbug(category, status, componentId)
	if err != nil {
		log.Error().Err(err).Send()
		return nil, InternalServerError(errInternalServerError)
	}

	return ResponseMsg{
		Message: "bugs are retrieved successfully",
		Data:    bugs,
	}, Ok()
}

// getBug retrieves a bug by its ID
// @Summary Get a bug
// @Description Retrieves a bug by its ID
// @Tags bugs
// @Produce json
// @Param id path string true "Bug ID"
// @Security ApiKeyAuth
// @Success 200 {object} ResponseMsg "bug is retrieved successfully (Bug details in the 'Data' field)"
// @Failure 400 {object} Response "Bad Request"
// @Failure 401 {object} Response "Unauthorized"
// @Failure 404 {object} Response "Not Found"
// @Failure 500 {object} Response "Internal Server Error"
// @Router /bug/{id} [get]
func (app *App) getBug(ctx *gin.Context) (interface{}, Response) {

	userId, exists := ctx.Get("user_id")
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

	members, err := app.getMemberstoAccessBug(bug)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Error().Err(err).Send()
		return nil, NotFound(errors.New("failed to find members that have access to this bug"))
	}

	if err != nil {
		log.Error().Err(err).Send()
		return nil, InternalServerError(errInternalServerError)
	}

	memberInProject := false

	// loop over the members to check admin state
	for _, member := range members {

		if userId == member.ID {

			memberInProject = true
			break

		}
	}

	if userId != bug.UserID && !memberInProject {
		return nil, Forbidden(errors.New("you have no access to get the bug"))
	}

	return ResponseMsg{
		Message: "bug is retrieved successfully",
		Data:    bug,
	}, Ok()
}

// updateBug updates a bug with the provided data
// @Summary Update a bug
// @Description Updates a bug with the provided data
// @Tags bugs
// @Accept json
// @Produce json
// @Param id path string true "Bug ID"
// @Param bug body updateBugInput true "Bug data"
// @Security ApiKeyAuth
// @Success 200 {object} ResponseMsg "bug is updated successfully"
// @Failure 400 {object} Response "Bad Request"
// @Failure 401 {object} Response "Unauthorized"
// @Failure 404 {object} Response "Not Found"
// @Failure 500 {object} Response "Internal Server Error"
// @Router /bug/{id} [put]
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

	userId, exists := ctx.Get("user_id")
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
	if userId != bug.UserID {
		return nil, Forbidden(errors.New("you have no access to update the bug"))
	}

	// proceed to update bug
	updatedBug := models.Bug{Summary: input.Summary, Category: input.Category, Severity: input.Severity, Status: input.Status, UpdatedAt: time.Now()}

	err = app.DB.UpdateBug(id, updatedBug)
	if err != nil {
		log.Error().Err(err).Send()
		return nil, InternalServerError(errInternalServerError)
	}

	return ResponseMsg{
		Message: "bug is updated successfully",
	}, Ok()
}

// deleteBug deletes a bug by its ID
// @Summary Delete a bug
// @Description Deletes a bug by its ID
// @Tags bugs
// @Produce json
// @Param id path string true "Bug ID"
// @Security ApiKeyAuth
// @Success 200 {object} ResponseMsg "bug is deleted successfully"
// @Failure 400 {object} Response "Bad Request"
// @Failure 401 {object} Response "Unauthorized"
// @Failure 404 {object} Response "Not Found"
// @Failure 500 {object} Response "Internal Server Error"
// @Router /bug/{id} [delete]
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

	members, err := app.getMemberstoAccessBug(bug)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Error().Err(err).Send()
		return nil, NotFound(errors.New("failed to find members that have access to this bug"))
	}

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

	if userId != bug.UserID && !adminMember {
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

// getMemberstoAccessBug a helper method that uses other methods to get members allowed to access a bug
func (app *App) getMemberstoAccessBug(bug models.Bug) ([]models.Member, error) {

	component, err := app.DB.GetComponent(strconv.Itoa(bug.ComponentID))

	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Error().Err(err).Send()
		return nil, errors.New("component is not found")
	}
	if err != nil {
		log.Error().Err(err).Send()
		return nil, errors.New("failed to get component due to internal server error")
	}

	project, err := app.DB.GetProject(component.ProjectID)

	if err != nil {
		log.Error().Err(err).Send()
		return nil, errors.New("failed to get project due to internal server error")
	}

	members, err := app.DB.GetMembersInProject(int(project.ID))
	if err != nil {
		log.Error().Err(err).Send()
		return nil, errors.New("failed to get members in the project due to internal server error")
	}

	return members, nil

}
