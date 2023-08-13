package app

import (
	"errors"
	"strconv"
	"time"

	//"log"

	"gorm.io/gorm"

	"github.com/codescalersinternships/Flyspray/models"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type CreateCommentInput struct {
	BugID   uint   `json:"bug_id" validate:"required"`
	Summary string `json:"summary" validate:"required"`
}

type updateCommentInput struct {
	Summary string `json:"summary" validate:"required"`
}

func (app *App) createComment(c *gin.Context) (interface{}, Response) {

	commentInput := CreateCommentInput{}

	if err := c.ShouldBindJSON(&commentInput); err != nil {
		log.Error().Err(err).Send()
		return nil, BadRequest(errors.New("bug ID, user ID, and summary of the comment should be provided"))
	}

	// user id will be taken from middlewares
	newComment := models.Comment{BugID: commentInput.BugID, Summary: commentInput.Summary, UserID: "1000", CreatedAt: time.Now()}

	if err := newComment.Validate(); err != nil {
		log.Error().Err(err).Send()
		return nil, BadRequest(errors.New("input data is invalid"))
	}

	comment, err := app.db.CreateComment(newComment)
	if err != nil {
		log.Error().Err(err).Send()
		return nil, InternalServerError(errors.New("failed to create comment"))

	}
	return ResponseMsg{
		Message: "Comment is created successfully",
		Data:    comment,
	}, Created()

}

func (app *App) getComment(c *gin.Context) (interface{}, Response) {

	idStr := c.Param("id")

	if idStr == "" {
		return nil, BadRequest(errors.New("comment ID is required"))
	}

	id, _ := strconv.ParseUint(idStr, 10, 64)

	comment, err := app.db.GetComment(uint(id))
	if err != nil {
		log.Error().Err(err).Send()
		return nil, NotFound(errors.New("comment is not found"))

	}

	return ResponseMsg{
		Message: "comment is found successfully",
		Data:    comment,
	}, Ok()

}

func (app *App) deleteComment(c *gin.Context) (interface{}, Response) {

	idStr := c.Param("id")

	if idStr == "" {
		return nil, BadRequest(errors.New("comment ID is required"))
	}

	id, _ := strconv.ParseUint(idStr, 10, 64)

	err := app.db.DeleteComment(uint(id))

	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Error().Err(err).Send()
		return nil, NotFound(errors.New("comment is not found"))

	}

	if err != nil {
		log.Error().Err(err).Send()
		return nil, InternalServerError(errors.New("failed to delete comment"))
	}
	return ResponseMsg{
		Message: "comment is deleted successfully",
	}, Ok()
}

func (app *App) listComments(c *gin.Context) (interface{}, Response) {

	bugIDStr := c.Query("bug_id")
	UserID := c.Query("user_id")

	bugID, err := strconv.ParseUint(bugIDStr, 10, 64)
	comments := app.db.ListComments(uint(bugID), UserID)

	if len(comments) == 0 {
		log.Error().Err(err).Send()
		return nil, NotFound(errors.New("no comments are found"))
	}

	return ResponseMsg{
		Message: "projects is retrieved successfully",
		Data:    comments,
	}, Ok()
}

func (app *App) updateComment(c *gin.Context) (interface{}, Response) {

	comment := updateCommentInput{}

	if err := c.ShouldBindJSON(&comment); err != nil {
		log.Error().Err(err).Send()
		return nil, BadRequest(errors.New("invalid comment data"))
	}

	idStr := c.Param("id")

	id, _ := strconv.ParseUint(idStr, 10, 64)

	err := app.db.UpdateComment(uint(id), comment.Summary)

	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Error().Err(err).Send()
		return nil, NotFound(errors.New("comment is not found"))
	}

	if err != nil {
		log.Error().Err(err).Send()
		return nil, InternalServerError(errors.New("failed to update the comment"))
	}

	return ResponseMsg{
		Message: "comment is updated successfully",
	}, Ok()
}
