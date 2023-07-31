package models

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DBClient used to start and make queries on database
type DBClient struct {
	Client *gorm.DB
}

// Start connects sqlite database
func (d *DBClient) Start(dbFilePath string) error {
	var err error
	d.Client, err = gorm.Open(sqlite.Open(dbFilePath), &gorm.Config{})

	return err
}
