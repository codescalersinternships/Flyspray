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

var ErrAccessDenied = errors.New("error access denied")

// Validate validates the comment struct using the validate tag
func (member *Member) Validate() error {
	validate := validator.New()
	return validate.Struct(member)
}

// CreateNewMember adds a new member to member table
func (db *DBClient) CreateNewMember(member Member, userId string) error {
	var authMember Member
	rows := db.Client.Model(&Member{}).Where("user_id = ? AND project_id = ?", userId, member.ProjectID).First(&authMember)
	if !authMember.Admin {
		var project Project
		if err := db.Client.Model(&Project{}).Where("id = ?", member.ProjectID).First(&project).Error; err != nil {
			return err
		}
		if project.OwnerID != userId {
			return ErrAccessDenied
		}
	}
	if rows.Error != nil && !errors.Is(rows.Error, gorm.ErrRecordNotFound) {
		return rows.Error
	}
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

// GetAllMembers returns all members in member table
func (db *DBClient) GetAllMembers() ([]Member, error) {
	var members []Member
	rows := db.Client.Find(&members)
	return members, rows.Error
}

// GetMembersInProject returns all member in a specific project
func (db *DBClient) GetMembersInProject(project_id int) ([]Member, error) {
	var members []Member
	rows := db.Client.Model(&Member{}).Where("project_id = ?", project_id).Find(&members)
	if rows.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return members, rows.Error
}

// UpdateMemberOwnership updates the admin bool in member table
func (db *DBClient) UpdateMemberOwnership(id int, admin bool, userId string) error {
	var member Member
	res := db.Client.Model(&member).Where("ID = ?", id).First(&member).Update("Admin", admin).Error
	var authMember Member
	rows := db.Client.Model(&Member{}).Where("user_id = ? AND project_id = ?", userId, member.ProjectID).First(&authMember)
	if !authMember.Admin {
		var project Project
		if err := db.Client.Model(&Project{}).Where("id = ?", member.ProjectID).First(&project).Error; err != nil {
			return err
		}
		if project.OwnerID != userId {
			return ErrAccessDenied
		}
	}
	if rows.Error != nil && !errors.Is(rows.Error, gorm.ErrRecordNotFound) {
		return rows.Error
	}
	return res
}
