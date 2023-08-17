package models

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// Member struct has Member table's columns
type Member struct {
	ID        int    `json:"id"`
	UserID    string `json:"user_id"  gorm:"not null;" validate:"required"`
	ProjectID int    `json:"project_id"  gorm:"not null" validate:"required"`
	Admin     bool   `json:"admin"`
}

// ErrAccessDenied indicates that user does not have access to create or update member
var ErrAccessDenied = errors.New("access denied	")

// Validate validates the comment struct using the validate tag
func (member *Member) Validate() error {
	validate := validator.New()
	return validate.Struct(member)
}

// CreateNewMember adds a new member to member table
func (db *DBClient) CreateNewMember(member Member) error {
	err := db.Client.Where("user_id = ? AND project_id = ?", member.UserID, member.ProjectID).First(&member).Error
	if err == nil {
		return gorm.ErrDuplicatedKey
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	res := db.Client.Create(&member)
	return res.Error
}

// GetMembersInProject returns all member in a specific project
func (db *DBClient) GetMembersInProject(project_id int) ([]Member, error) {
	var members []Member
	rows := db.Client.Model(&Member{}).Where("project_id = ?", project_id).Find(&members)
	return members, rows.Error
}

// UpdateMemberOwnership updates the admin bool in member table
func (db *DBClient) UpdateMemberOwnership(id int, admin bool, userId string) error {
	var member Member
	res := db.Client.Model(&member).Where("ID = ?", id).First(&member)
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	var authMember Member
	rows := db.Client.Model(&Member{}).Where("user_id = ? AND project_id = ?", userId, member.ProjectID).First(&authMember)
	if !authMember.Admin {
		return ErrAccessDenied
	}
	if rows.Error != nil && !errors.Is(rows.Error, gorm.ErrRecordNotFound) {
		return rows.Error
	}
	resErr := res.Update("Admin", admin).Error
	return resErr
}

func (db *DBClient) CheckUserAccess(member Member, userId string) error {
	var authMember Member
	rows := db.Client.Model(&Member{}).Where("user_id = ? AND project_id = ?", userId, member.ProjectID).First(&authMember)

	if !authMember.Admin {
		return ErrAccessDenied
	}

	if rows.Error != nil && !errors.Is(rows.Error, gorm.ErrRecordNotFound) {
		return rows.Error
	}
	return nil
}
