package models_test

import (
	"testing"
	"time"

	"github.com/codescalersinternships/Flyspray/models"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestValidateComment(t *testing.T) {

	comment := models.Comment{
		ID:        1,
		UserID:    "2",
		BugID:     3,
		Summary:   "This is a comment",
		CreatedAt: time.Now(),
	}

	validate := validator.New()
	err := validate.Struct(comment)
	assert.NoError(t, err, "comment should be valid")

	comment.UserID = ""
	err = validate.Struct(comment)
	assert.Error(t, err, "comment should be invalid")
}
