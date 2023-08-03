package app

import (
	"errors"
	"log"
	"net/http"

	"gorm.io/gorm"

	"github.com/codescalersinternships/Flyspray/models"
	"github.com/gin-gonic/gin"
)

// Response is a struct that holds the ok response
type Response struct {
	Message string           `json:"message"`
	Payload []models.Comment `json:"payload"`
}

// ErrorResponse is a struct that holds the error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// CreateComment creates a comment on a specific bug
func (app *App) CreateComment(c *gin.Context) {

	comment := models.Comment{}

	if err := c.ShouldBindJSON(&comment); err != nil {
		log.Println("failed to create the comment:", err)
		c.JSON(http.StatusBadRequest,
			ErrorResponse{Error: "bugID,ownerID and summary of the comment should be provided "})
		return
	}

	if err := comment.Validate(); err != nil {
		log.Println("bugID,ownerID and summary of the comment should be validated", err)
		c.JSON(http.StatusBadRequest,
			ErrorResponse{Error: "comment is not validated"})
		return
	}

	comment, err := app.client.CreateComment(comment)
	if err != nil {
		log.Fatal(err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to create comment"})
		return
	}

	c.JSON(http.StatusCreated,
		Response{Message: "comment is created",
			Payload: []models.Comment{comment}})

}

// GetComment gets a comment on a specific bug by id
func (app *App) GetComment(c *gin.Context) {

	id := c.Param("id")

	if id == "" {
		log.Println("comment ID is required")
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "comment ID is required"})
		return
	}

	comment, err := app.client.GetComment(id)
	if err != nil {
		log.Println("comment is not found:", err)
		c.JSON(http.StatusNotFound,
			ErrorResponse{Error: "comment is not found"})
		return

	}

	c.JSON(http.StatusOK,
		Response{Message: "comment is found",
			Payload: []models.Comment{comment}})

}

// DeleteComment deletes a comment on a specific bug by id
func (app *App) DeleteComment(c *gin.Context) {

	id := c.Param("id")

	if id == "" {
		log.Println("comment ID is required")
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "comment ID is required"})
		return
	}

	if err := app.client.DeleteComment(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("comment is not found:", err)
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "comment is not found"})
		} else {
			log.Fatal(err)
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to delete comment"})
		}
		return
	}

	c.JSON(http.StatusOK, Response{Message: "comment is deleted"})
}

// ListComments lists all the comments for a specific bug
func (app *App) ListComments(c *gin.Context) {

	bugID := c.Query("bug_id")

	if bugID == "" {
		log.Println("bug id is not found")
		c.JSON(http.StatusBadRequest,
			ErrorResponse{Error: "bug id is not found"})
		return
	}

	comments := app.client.ListComments(bugID)

	if len(comments) == 0 {
		log.Println("no comments are found for this bug")
		c.JSON(http.StatusNotFound,
			ErrorResponse{Error: "no comments are found for this bug"})
		return
	}

	c.JSON(http.StatusOK,
		Response{Message: "comment are found for the bug",
			Payload: comments})
}

// UpdateComment updates a comment on a specific bug by id
func (app *App) UpdateComment(c *gin.Context) {

	comment := models.Comment{}
	id := c.Param("id")

	if err := c.ShouldBindJSON(&comment); err != nil {
		log.Println("failed to update the comment:", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid comment data"})
		return
	}

	if err := comment.Validate(); err != nil {
		log.Println("bugID, ownerID, and summary of the comment should be validated", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "comment is not validated"})
		return
	}

	comment, err := app.client.UpdateComment(comment, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("comment is not found:", err)
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "comment is not found"})
		} else {
			log.Fatal(err)
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to update comment"})
		}
		return
	}

	c.JSON(http.StatusOK, Response{
		Message: "comment is updated",
		Payload: []models.Comment{comment},
	})
}
