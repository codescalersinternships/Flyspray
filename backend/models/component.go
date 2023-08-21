package models

import (
	"time"

	"gorm.io/gorm"
)

// Component is a struct representing a Component
type Component struct {
	// ID uniqe for each Component
	ID uint `json:"id" gorm:"primaryKey;autoIncrement:true"`
	// UserID descripe user create copmonent
	UserID string `json:"user_id" validate:"required"`
	// ProjectID descripe id project
	ProjectID string `json:"project_id" validate:"required"`
	// Name descripe each Component
	Name string `json:"name" gorm:"unique" validate:"required "`
	// CreatedAt descripe time for the component
	CreatedAt time.Time `json:"created_at"`
}

// CreateComponent adds new component to database
func (d *DBClient) CreateComponent(c Component) (Component, error) {
	return c, d.Client.Create(&c).Error
}

// UpdateComponent updates component
func (d *DBClient) UpdateComponent(id string, updatedComponent Component) error {
	result := d.Client.Model(&updatedComponent).Where("id = ?", id).
		Update("name", updatedComponent.Name).
		Update("project_id", updatedComponent.ProjectID)

	if result.Error == nil && result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}

// GetComponent gets compoennt by id
func (d *DBClient) GetComponent(id string) (Component, error) {
	c := Component{}
	return c, d.Client.First(&c, id).Error
}

// FilterComponents filters all components by project id, component name, creation date
func (d *DBClient) FilterComponents(ProjectID, componentName, date string) ([]Component, error) {
	projects := []Component{}

	query := d.Client
	if ProjectID != "" {
		query = query.Where("project_id = ?", ProjectID)
	}
	if componentName != "" {
		query = query.Where("name = ?", componentName)
	}
	if date != "" {
		query = query.Where("created_at > ?", date)
	}

	return projects, query.Find(&projects).Error
}

// DeleteComponent deletes component by id
func (d *DBClient) DeleteComponent(id string) error {
	result := d.Client.Delete(&Component{}, id)
	if result.Error == nil && result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}
