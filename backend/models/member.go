package models

import (
	_ "embed"
)

// Member struct has Member table's columns
type Member struct {
	ID        int  `json:"id"`
	UserID    int  `json:"user_id"`
	ProjectID int  `json:"project_id"`
	Admin     bool `json:"admin_bool"`
}

// Crete adds a new member to member table
func (db *DBClient) CreateNewMember(member Member) error {
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
	return res.Error
}
