package models

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// Member struct has Member table's columns
type Member struct {
	ID        int       `json:"id"`
	UserID    string    `json:"user_id"  gorm:"not null;" validate:"required"`
	ProjectID int       `json:"project_id"  gorm:"not null" validate:"required"`
	Admin     bool      `json:"admin_bool"`
	Project   []Project `gorm:"many2many:member_project;"`
}

// Validate validates the comment struct using the validate tag
func (member *Member) Validate() error {
	validate := validator.New()
	return validate.Struct(member)
}

// CreateNewMember adds a new member to member table
func (db *DBClient) CreateNewMember(member Member) (*Member, error) {
	err := db.Client.Where("user_id = ? AND project_id = ?", member.UserID, member.ProjectID).First(&member).Error
	if err == nil {
		return nil, gorm.ErrDuplicatedKey
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	res := db.Client.Create(&member)
	if res.Error != nil {
		return nil, res.Error
	}
	// Retrieve the newly created member from the database
	err = db.Client.First(&member, member.ID).Error
	if err != nil {
		return nil, err
	}
	return &member, nil
}

// GetAllMembers returns all members in member table
func (db *DBClient) GetAllMembers() ([]Member, error) {
	var members []Member
	rows := db.Client.Find(&members)
	return members, rows.Error
}

// UpdateMemberOwnership updates the admin bool in member table
func (db *DBClient) UpdateMemberOwnership(member Member, id int) (*Member, error) {
	res := db.Client.Model(&Member{}).Where("ID = ?", id).Updates(map[string]interface{}{
		"Admin": member.Admin,
	})
	if res.Error != nil {
		return nil, res.Error
	}

	if res.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	updatedMember := Member{}
	if err := db.Client.First(&updatedMember, id).Error; err != nil {
		return nil, err
	}

	return &updatedMember, nil
}
