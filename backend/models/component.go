package models

import (
	"time"
)

// Component is a struct representing a Component
type Component struct {
	// ID uniqe for each Component
	ID uint `json:"id"`
	// ProjectID descripe id project
	ProjectID uint `json:"project_id" validate:"required"`
	// Name descripe each Component
	Name string `json:"name" validate:"required"`
	// CreatedAt descripe time for the component
	CreatedAt time.Time `json:"created_at"`
}
