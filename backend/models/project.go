package models

import (
	"time"
)

type Project struct {
	ID        uint      `json:"id" gorm:"primary_key; unique"`
	Name      string    `json:"name" validate:"nonzero"`
	OwnerId   uint      `json:"owner_id" validate:"nonzero"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateProject adds new project to database
func (d *DBClient) CreateProject(project Project) (Project, error) {
	result := d.Client.Create(&project)
	return project, result.Error
}

// UpdateProject updates project
func (d *DBClient) UpdateProject(id string, updatedProject Project) (Project, error) {
	project := Project{}
	result := d.Client.First(&project, id)
	if result.Error != nil {
		return Project{}, result.Error
	}

	project.Name = updatedProject.Name
	d.Client.Save(&project)

	return project, nil
}

// GetProject gets project by id
func (d *DBClient) GetProject(id string) (Project, error) {
	project := Project{}
	result := d.Client.First(&project, id)
	return project, result.Error
}

// FilterProjects filters all projects by user id, project name, creation date
func (d *DBClient) FilterProjects(userId, projectName, date string) []Project {
	projects := []Project{}

	if userId == "" && projectName == "" && date == "" {
		d.Client.Find(&projects)
		return projects
	}

	if userId == "" && projectName == "" {
		d.Client.Where("created_at > ?", date).Find(&projects)
		return projects
	}
	if userId == "" && date == "" {
		d.Client.Where("name = ?", projectName).Find(&projects)
		return projects
	}
	if projectName == "" && date == "" {
		d.Client.Where("owner_id = ?", userId).Find(&projects)
		return projects
	}

	if userId == "" {
		d.Client.Where("name = ? and created_at > ?", projectName, date).Find(&projects)
		return projects
	}
	if projectName == "" {
		d.Client.Where("owner_id = ? and created_at > ?", userId, date).Find(&projects)
		return projects
	}
	if date == "" {
		d.Client.Where("owner_id = ? and name = ?", userId, projectName).Find(&projects)
		return projects
	}

	d.Client.Where("owner_id = ? and name = ? and created_at > ?", userId, projectName, date).Find(&projects)
	return projects
}

// DeleteProject deletes project by id
func (d *DBClient) DeleteProject(id string) {
	project := Project{}
	d.Client.Delete(&project, id)
}
