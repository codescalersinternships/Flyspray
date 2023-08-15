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

func (app *App) getComment(c *gin.Context) (interface{}, Response) {

	idStr := c.Param("id")

	if idStr == "" {
		return nil, BadRequest(errors.New("comment id is required"))
	}

	id, _ := strconv.ParseUint(idStr, 10, 64)

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

func (app *App) deleteComment(c *gin.Context) (interface{}, Response) {

	idStr := c.Param("id")

	if idStr == "" {
		return nil, BadRequest(errors.New("comment id is required"))
	}

	id, _ := strconv.ParseUint(idStr, 10, 64)

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

	err = app.DB.DeleteComment(uint(id))

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
	return ResponseMsg{
		Message: "comment is deleted successfully",
	}, Ok()
}

func (app *App) listComments(c *gin.Context) (interface{}, Response) {

	bugIDStr := c.Query("bug_id")
	UserID := c.Query("user_id")

	bugID, _ := strconv.ParseUint(bugIDStr, 10, 64)
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

func (app *App) updateComment(c *gin.Context) (interface{}, Response) {

	var updatedComment updateCommentInput

	if err := c.ShouldBindJSON(&updatedComment); err != nil {
		log.Error().Err(err).Send()
		return nil, BadRequest(errors.New("invalid comment data"))
	}

	idStr := c.Param("id")

	id, _ := strconv.ParseUint(idStr, 10, 64)

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
