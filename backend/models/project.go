package models

import (
	"time"

	"gorm.io/gorm"
)

// Project model
type Project struct {
	ID        uint      `json:"id" gorm:"primary_key; unique"`
	Name      string    `json:"name" validate:"nonzero" gorm:"unique"`
	OwnerId   uint      `json:"owner_id" validate:"nonzero"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateProject adds new project to database
func (d *DBClient) CreateProject(p Project) (Project, error) {
	return p, d.Client.Create(&p).Error
}

// UpdateProject updates project
func (d *DBClient) UpdateProject(id string, updatedProject Project) error {
	return d.Client.Model(&updatedProject).Where("id = ?", id).Update("name", updatedProject.Name).Error
}

// GetProject gets project by id
func (d *DBClient) GetProject(id string) (Project, error) {
	p := Project{}
	return p, d.Client.First(&p, id).Error
}

// GetProjectByName gets project by name
func (d *DBClient) GetProjectByName(name string) (Project, error) {
	p := Project{}
	return p, d.Client.Where("name = ?", name).First(&p).Error
}

// FilterProjects filters all projects by user id, project name, creation date
func (d *DBClient) FilterProjects(ownerId, projectName, date string) ([]Project, error) {
	projects := []Project{}

	query := d.Client
	if ownerId != "" {
		query = query.Where("owner_id = ?", ownerId)
	}
	if projectName != "" {
		query = query.Where("name = ?", projectName)
	}
	if date != "" {
		query = query.Where("created_at > ?", date)
	}

	return projects, query.Find(&projects).Error
}

// DeleteProject deletes project by id
func (d *DBClient) DeleteProject(id string) error {
	result := d.Client.Delete(&Project{}, id)
	if result.Error == nil && result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return result.Error
}
