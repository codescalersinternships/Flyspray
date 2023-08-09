package models

import (
	"time"

	"gorm.io/gorm"
)

type Bug struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement:true"`
	UserID      string    `json:"user_id"`
	ComponentID int       `json:"component_id"`
	Category    string    `json:"category"`
	Severity    string    `json:"severity"`
	Summary     string    `json:"summary"`
	Status      string    `json:"status"`
	Votes       int       `json:"votes"`
	Opened      bool      `json:"opened"`
	OpenedAt    time.Time `json:"opened_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateNewBug create new Bug to database
func (d *DBClient) CreateNewBug(p Bug) (Bug, error) {
	return p, d.Client.Create(&p).Error
}

// GetSpecificBug get bug using id from database
func (d *DBClient) GetSpecificBug(id string) (Bug, error) {
	p := Bug{}
	return p, d.Client.Where("id = ?", id).First(&p).Error
}

func (d *DBClient) UpdateBug(id string, updateBug Bug) *gorm.DB {
	return d.Client.Model(&Bug{}).Where("id = ?", id).Updates(updateBug)
}

// FilterBugs filters all bugs by user id, bug category, bug status, component id, opened,
func (d *DBClient) FilterBugs(userId, category, status, component_id string) ([]Bug, error) {
	projects := []Bug{}

	query := d.Client

	if userId != "" {
		query = query.Where("user_id = ?", userId)
	}

	if category != "" {
		query = query.Where("category = ?", category)
	}

	if component_id != "" {
		query = query.Where("component_id = ?", component_id)
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	return projects, query.Find(&projects).Error
}

// DeleteBug delete bug from database
func (d *DBClient) DeleteBug(id string) error {
	return d.Client.Delete(&Bug{}, id).Error
}
