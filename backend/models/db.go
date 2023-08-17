package models

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// NewDBClient connects sqlite database and returns DBClient
func NewDBClient(dbFilePath string) (DBClient, error) {
	d := DBClient{}
	// connect database
	var err error
	d.Client, err = gorm.Open(sqlite.Open(dbFilePath), &gorm.Config{})

	return d, err
}

// DBClient used to start and make queries on database
type DBClient struct {
	Client *gorm.DB
}

// Migrate makes migrations for the database
func (d *DBClient) Migrate() error {
	return d.Client.AutoMigrate(&Project{}, &Comment{}, &Member{})
}
