package models

import (
	"time"

	"github.com/go-playground/validator/v10"
)

// Comment is a struct representing comments on bugs
type Comment struct {

	// ID represents the unique id of each comment
	ID uint `json:"id"`

	// OwnerID represents the unique id of the user that wrote the comment
	OwnerID uint `json:"owner_id" validate:"required"`

	// BugID represents the unique id of the bug
	BugID uint `json:"bug_id" validate:"required"`

	// Summary represents the text of the comment
	Summary string `json:"summary" validate:"required"`

	// CreatedAt represents the time at which the comment was created
	CreatedAt time.Time
}

// Validate validates the comment struct using the validate tag
func (c *Comment) Validate() error {
	validate := validator.New()
	return validate.Struct(c)
}
