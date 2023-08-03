package models

import (
	"time"

	"github.com/go-playground/validator/v10"
)

// Comment is a struct representing comments on bugs
type Comment struct {

	// ID represents the unique id of each comment
	ID uint `gorm:"primaryKey" json:"id"`

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

// CreateComment creates a comment on a specific bug
func (db *DBClient) CreateComment(comment Comment) (Comment, error) {

	result := db.Client.Create(&comment)
	return comment, result.Error
}

// GetComment gets a comment on a specific bug by id
func (db *DBClient) GetComment(id string) (Comment, error) {

	comment := Comment{}

	result := db.Client.First(&comment, id)
	return comment, result.Error

}

// DeleteComment deletes a comment on a specific bug by id
func (db *DBClient) DeleteComment(id string) error {

	comment := Comment{}

	if result := db.Client.First(&comment, id); result.Error != nil {
		return result.Error
	}

	if result := db.Client.Delete(&comment, id); result.Error != nil {
		return result.Error
	}

	return nil

}

// ListComments lists all the comments for a specific bug
func (db *DBClient) ListComments(bugID string) []Comment {

	comments := []Comment{}
	db.Client.Where("bug_id=?", bugID).Find(&comments)

	return comments

}

// UpdateComment updates a comment on a specific bug by id
func (db *DBClient) UpdateComment(comment Comment, id string) (Comment, error) {

	if result := db.Client.First(&comment, id); result.Error != nil {
		return comment, result.Error
	}

	result := db.Client.Save(&comment)
	return comment, result.Error

}
