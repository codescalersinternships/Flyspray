package models

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

// ErrNoAccess is returned when a user not allowed to get certain bugs tries to get them
var ErrNoAccess = errors.New("user has no access to get bugs")

type Bug struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement:true"`
	UserID      string    `json:"user_id" validate:"required"`
	ComponentID int       `json:"component_id" validate:"required"`
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
func (d *DBClient) GetBug(id string) (Bug, error) {
	p := Bug{}
	return p, d.Client.Where("id = ?", id).First(&p).Error
}

func (d *DBClient) UpdateBug(id string, updatedBug Bug) error {
	result := d.Client.Model(&updatedBug).Where("id = ?", id).Updates(updatedBug)
	if result.Error == nil && result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}

// Filterbug filters all bugs by user id, bug category, bug status, component id
func (d *DBClient) Filterbug(category, status, component_id string) ([]Bug, error) {
	bug := []Bug{}

	query := d.Client

	if category != "" {
		query = query.Where("category = ?", category)
	}

	if component_id != "" {
		query = query.Where("component_id = ?", component_id)
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}
	return bug, query.Find(&bug).Error
}

// DeleteBug delete bug from database
func (d *DBClient) DeleteBug(id string) error {
	bug := Bug{}

	result := d.Client.Delete(&bug, id)

	if result.Error == nil && result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return result.Error
}
