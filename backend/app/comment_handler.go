package app

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"gorm.io/gorm"

	"github.com/codescalersinternships/Flyspray/models"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type createCommentInput struct {
	BugID   uint   `json:"bug_id" validate:"required"`
	Summary string `json:"summary" validate:"required"`
}

type updateCommentInput struct {
	Summary string `json:"summary" validate:"required"`
}

// createComment creates a new comment for a bug
// @Summary Creates a comment
// @Description Creates a new comment for a bug in the database
// @Tags comments
// @Accept json
// @Produce json
// @Param commentInput body createCommentInput true "Comment input object"
// @Security ApiKeyAuth
// @Success 201 {object} ResponseMsg "Comment is created successfully (Comment details in the 'Data' field)"
// @Failure 400 {object} Response "Bad Request"
// @Failure 401 {object} Response "Unauthorized"
// @Failure 500 {object} Response "Internal Server Error"
// @Router /comment [post]
func (app *App) createComment(c *gin.Context) (interface{}, Response) {

	var commentInput createCommentInput

	if err := c.ShouldBindJSON(&commentInput); err != nil {
		log.Error().Err(err).Send()
		return nil, BadRequest(errors.New("bug id, user id, and summary of the comment should be provided"))
	}

	userID, exists := c.Get("user_id")
	if !exists {
		return nil, UnAuthorized(errors.New("authentication is required"))
	}

	newComment := models.Comment{BugID: commentInput.BugID, Summary: commentInput.Summary, UserID: userID.(string), CreatedAt: time.Now()}

	if err := newComment.Validate(); err != nil {
		log.Error().Err(err).Send()
		return nil, BadRequest(errors.New("input data is invalid"))
	}

	comment, err := app.DB.CreateComment(newComment)
	if err != nil {
		log.Error().Err(err).Send()
		return nil, InternalServerError(errInternalServerError)

	}
	return ResponseMsg{
		Message: "Comment is created successfully",
		Data:    comment,
	}, Created()

}

// getComment retrieves a comment by ID
// @Summary Retrieves a comment
// @Description Retrieves a comment by its ID from the database
// @Tags comments
// @Produce json
// @Param id path int true "Comment ID"
// @Security ApiKeyAuth
// @Success 200 {object} ResponseMsg "comment is found successfully (Comment details in the 'Data' field)"
// @Failure 400 {object} Response "Bad Request"
// @Failure 401 {object} Response "Unauthorized"
// @Failure 404 {object} Response "Not Found"
// @Failure 500 {object} Response "Internal Server Error"
// @Router /comment/{id} [get]
func (app *App) getComment(c *gin.Context) (interface{}, Response) {

	idStr := c.Param("id")

	if idStr == "" {
		return nil, BadRequest(errors.New("comment id is required"))
	}

	id, err := strconv.ParseUint(idStr, 10, 64)

	if err != nil {
		return nil, BadRequest(errors.New("failed to parse the comment id"))
	}

	comment, err := app.DB.GetComment(uint(id))

	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Error().Err(err).Send()
		return nil, NotFound(errors.New("comment is not found"))

	}

	if err != nil {
		log.Error().Err(err).Send()
		return nil, InternalServerError(errInternalServerError)
	}

	return ResponseMsg{
		Message: "comment is found successfully",
		Data:    comment,
	}, Ok()

}

// deleteComment deletes a comment by ID
// @Summary Deletes a comment
// @Description Deletes a comment by its ID from the database
// @Tags comments
// @Produce json
// @Param id path int true "Comment ID"
// @Security ApiKeyAuth
// @Success 200 {object} ResponseMsg "comment is deleted successfully"
// @Failure 400 {object} Response "Bad Request"
// @Failure 401 {object} Response "Unauthorized"
// @Failure 403 {object} Response "Forbidden"
// @Failure 404 {object} Response "Not Found"
// @Failure 500 {object} Response "Internal Server Error"
// @Router /comment/{id} [delete]
func (app *App) deleteComment(c *gin.Context) (interface{}, Response) {

	idStr := c.Param("id")

	if idStr == "" {
		return nil, BadRequest(errors.New("comment id is required"))
	}

	id, err := strconv.ParseUint(idStr, 10, 64)

	if err != nil {
		return nil, BadRequest(errors.New("failed to parse the comment id"))
	}

	userID, exists := c.Get("user_id")
	if !exists {
		return nil, UnAuthorized(errors.New("authentication is required"))
	}

	comment, err := app.DB.GetComment(uint(id))

	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Error().Err(err).Send()
		return nil, NotFound(errors.New("comment is not found"))

	}

	if err != nil {
		log.Error().Err(err).Send()
		return nil, InternalServerError(errInternalServerError)
	}

	if fmt.Sprint(userID) != fmt.Sprint(comment.UserID) {
		return nil, Forbidden(errors.New("you have no access to delete the comment"))
	}

	err = app.DB.DeleteComment(uint(id))

	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Error().Err(err).Send()
		return nil, NotFound(errors.New("comment is not found"))

	}

	if err != nil {
		log.Error().Err(err).Send()
		return nil, InternalServerError(errInternalServerError)
	}

	return ResponseMsg{
		Message: "comment is deleted successfully",
	}, Ok()
}

// listComments retrieves a list of comments
// @Summary Retrieves comments
// @Description Retrieves a list of comments from the database
// @Tags comments
// @Produce json
// @Param bug_id query int false "Bug ID"
// @Param user_id query string false "User ID"
// @Success 200 {object} ResponseMsg "comments are retrieved successfully (Comments details in the 'Data' field)"
// @Failure 400 {object} Response "Bad Request"
// @Failure 500 {object} Response "Internal Server Error"
// @Router /comment/filters [get]
func (app *App) listComments(c *gin.Context) (interface{}, Response) {

	bugIDStr := c.Query("bug_id")
	UserID := c.Query("user_id")

	bugID, err := strconv.ParseUint(bugIDStr, 10, 64)

	if bugIDStr != "" && err != nil {
		return nil, BadRequest(errors.New("failed to parse the bug id"))

	}

	comments := app.DB.ListComments(uint(bugID), UserID)

	if len(comments) == 0 {
		return ResponseMsg{
			Message: "no comments are found",
			Data:    comments,
		}, Ok()

	}

	return ResponseMsg{
		Message: "comments are retrieved successfully",
		Data:    comments,
	}, Ok()
}

// updateComment updates a comment by ID
// @Summary Updates a comment
// @Description Updates a comment by its ID in the database
// @Tags comments
// @Accept json
// @Produce json
// @Param id path int true "Comment ID"
// @Param input body updateCommentInput true "Updated comment data"
// @Security ApiKeyAuth
// @Success 200 {object} ResponseMsg "comment is updated successfully"
// @Failure 400 {object} Response "Bad Request"
// @Failure 401 {object} Response "Unauthorized"
// @Failure 403 {object} Response "Forbidden"
// @Failure 404 {object} Response "Not Found"
// @Failure 500 {object} Response "Internal Server Error"
// @Router /comment/{id} [put]
func (app *App) updateComment(c *gin.Context) (interface{}, Response) {

	var updatedComment updateCommentInput

	if err := c.ShouldBindJSON(&updatedComment); err != nil {
		log.Error().Err(err).Send()
		return nil, BadRequest(errors.New("invalid comment data"))
	}

	idStr := c.Param("id")

	if idStr == "" {
		return nil, BadRequest(errors.New("comment id is required"))
	}

	id, err := strconv.ParseUint(idStr, 10, 64)

	if err != nil {
		return nil, BadRequest(errors.New("failed to parse the comment id"))
	}

	userID, exists := c.Get("user_id")
	if !exists {
		return nil, UnAuthorized(errors.New("authentication is required"))
	}

	comment, err := app.DB.GetComment(uint(id))

	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Error().Err(err).Send()
		return nil, NotFound(errors.New("comment is not found"))

	}

	if err != nil {
		log.Error().Err(err).Send()
		return nil, InternalServerError(errInternalServerError)
	}

	if fmt.Sprint(userID) != fmt.Sprint(comment.UserID) {
		return nil, Forbidden(errors.New("you have no access to update the comment"))
	}

	err = app.DB.UpdateComment(uint(id), updatedComment.Summary)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Error().Err(err).Send()
		return nil, NotFound(errors.New("comment is not found"))
	}

	if err != nil {
		log.Error().Err(err).Send()
		return nil, InternalServerError(errInternalServerError)
	}

	return ResponseMsg{
		Message: "comment is updated successfully",
	}, Ok()
}
