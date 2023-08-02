package handlers

import (
	"log"
	"net/http"

	"github.com/codescalersinternships/Flyspray/models"
	"github.com/gin-gonic/gin"
)

// CreateComment creates a comment on a specific bug
func CreateComment(c *gin.Context, db models.DBClient) {

	comment := models.Comment{}

	if err := c.ShouldBindJSON(&comment); err != nil {
		log.Println("error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to create comment"})
		return
	}

	// validate the comment object
	if err := comment.Validate(); err != nil {
		log.Println("error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid comment data"})
		return
	}

	// create the comment object
	if result := db.Client.Create(&comment); result.Error != nil {
		log.Fatal(result.Error)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, comment)
}

// GetComment gets a comment on a specific bug by id
func GetComment(c *gin.Context, db models.DBClient) {

	comment := models.Comment{}

	id := c.Param("id")

	if result := db.Client.First(&comment, id); result.Error != nil {
		log.Println("error:", result.Error)
		c.JSON(http.StatusNotFound, gin.H{"error": "comment is not found"})
		return
	}
	c.JSON(http.StatusOK, comment)
}

// DeleteComment deletes a comment on a specific bug by id
func DeleteComment(c *gin.Context, db models.DBClient) {

	comment := models.Comment{}

	id := c.Param("id")

	if result := db.Client.First(&comment, id); result.Error != nil {
		log.Println("error:", result.Error)
		c.JSON(http.StatusNotFound, gin.H{"error": "comment is not found"})
		return
	}

	if result := db.Client.Delete(&comment, id); result.Error != nil {
		log.Fatal(result.Error)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "comment is deleted"})
}

// ListComments lists all the comments for a specific bug
func ListComments(c *gin.Context, db models.DBClient) {

	bugID := c.Query("bug_id")

	if bugID == "" {

		c.JSON(http.StatusBadRequest, gin.H{"error": "bug id is not found"})
		return
	}

	comments := []models.Comment{}
	db.Client.Where(" = ?", bugID).Find(&comments)

	if len(comments) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "no comments found for this bug"})
		return
	}

	c.JSON(http.StatusOK, comments)
}
