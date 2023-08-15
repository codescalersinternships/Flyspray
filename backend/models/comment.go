package models

import (
	"time"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// Comment is a struct representing comments on bugs
type Comment struct {

	// ID represents the unique id of each comment
	ID uint `gorm:"primaryKey" json:"id"`

	// UserID represents the unique id of the user that wrote the comment
	UserID string `json:"user_id" validate:"required"`

	// BugID represents the unique id of the bug
	BugID uint `json:"bug_id" validate:"required"`

	// Summary represents the text of the comment
	Summary string `json:"summary" validate:"required"`

	// CreatedAt represents the time at which the comment was created
	CreatedAt time.Time `json:"created_at"`
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
func (db *DBClient) GetComment(id uint) (Comment, error) {

	comment := Comment{}

	result := db.Client.First(&comment, id)

	return comment, result.Error

}

// DeleteComment deletes a comment on a specific bug by id
func (db *DBClient) DeleteComment(id uint) error {

	comment := Comment{}

	result := db.Client.Delete(&comment, id)

	if result.Error == nil && result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return result.Error

}

// ListComments lists all the comments for a specific bug
func (db *DBClient) ListComments(bugID uint, UserID string) []Comment {

	comments := []Comment{}

	if bugID != 0 || UserID != "" {
		db.Client.Where("bug_id = ? OR user_id = ?", bugID, UserID).Find(&comments)
	} else {
		db.Client.Find(&comments)
	}

	return comments
}

// UpdateComment updates a comment on a specific bug by id
func (db *DBClient) UpdateComment(id uint, newSummary string) error {

	return db.Client.Model(&Comment{}).Where("id = ?", id).Update("summary", newSummary).Error

}
