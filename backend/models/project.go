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
func CreateProject(project Project, db DBClient) (Project, error) {
	result := db.Client.Create(&project)
	return project, result.Error
}

// UpdateProject updates project
func UpdateProject(id string, updatedProject Project, db DBClient) (Project, error) {
	project := Project{}
	result := db.Client.First(&project, id)
	if result.Error != nil {
		return Project{}, result.Error
	}

	project.Name = updatedProject.Name
	db.Client.Save(&project)

	return project, nil
}

// GetProject gets project by id
func GetProject(id string, db DBClient) (Project, error) {
	project := Project{}
	result := db.Client.First(&project, id)
	return project, result.Error
}

// FilterProjects filters all projects by user id, project name, creation date
func FilterProjects(userId, projectName, date string, db DBClient) []Project {
	projects := []Project{}

	if userId == "" && projectName == "" && date == "" {
		db.Client.Find(&projects)
		return projects
	}

	if userId == "" && projectName == "" {
		db.Client.Where("created_at > ?", date).Find(&projects)
		return projects
	}
	if userId == "" && date == "" {
		db.Client.Where("name = ?", projectName).Find(&projects)
		return projects
	}
	if projectName == "" && date == "" {
		db.Client.Where("owner_id = ?", userId).Find(&projects)
		return projects
	}

	if userId == "" {
		db.Client.Where("name = ? and created_at > ?", projectName, date).Find(&projects)
		return projects
	}
	if projectName == "" {
		db.Client.Where("owner_id = ? and created_at > ?", userId, date).Find(&projects)
		return projects
	}
	if date == "" {
		db.Client.Where("owner_id = ? and name = ?", userId, projectName).Find(&projects)
		return projects
	}

	db.Client.Where("owner_id = ? and name = ? and created_at > ?", userId, projectName, date).Find(&projects)
	return projects
}

// DeleteProject deletes project by id
func DeleteProject(id string, db DBClient) {
	project := Project{}
	db.Client.Delete(&project, id)
}
