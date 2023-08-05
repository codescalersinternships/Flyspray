package models

import "errors"

// Member struct has Member table's columns
type Member struct {
	ID        int  `json:"id"`
	UserID    int  `json:"user_id"  gorm:"not null"`
	ProjectID int  `json:"project_id"  gorm:"not null"`
	Admin     bool `json:"admin_bool"`
}

var (
	ErrEmptyMemberFields = errors.New("userid and projectid cannot be null")
	ErrMemberNotFound    = errors.New("member was not found")
)

// Crete adds a new member to member table
func (db *DBClient) CreateNewMember(member Member) error {
	if member.UserID == 0 || member.ProjectID == 0 {
		return ErrEmptyMemberFields
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

// UpdateMemberOwnership updates the admin bool in member table
func (db *DBClient) UpdateMemberOwnership(member Member, id int) error {
	res := db.Client.Model(&Member{}).Where("ID = ?", id).Update("Admin", member.Admin)
	if res.RowsAffected == 0{
		return ErrMemberNotFound
	}
	return res.Error
}

