package models

import (
	"fmt"

	"gopkg.in/validator.v2"
)

type User struct {
	ID                string `gorm:"primary_key;<-:false"`
	Name              string `json:"name" validate:"required"`
	Email             string `json:"email" gorm:"unique;not null" validate:"required regexp=^[0-9a-z]+@[0-9a-z]+(\\.[0-9a-z]+)+$"`
	Password          string `gorm:"not null" validate:"required"`
	Verification_code string `gorm:"unique;not null"`
	Verified          uint   `gorm:"default:0"`
}

func (u *User) Validate() error {
	return validator.Validate(u)
}

func (db *DBClient) CreateUser(user User) (User, error) {

	result := db.Client.Create(&user)
	return user, result.Error
}

func (db *DBClient) UpdateUser(user User) {
	result := db.Client.First(&user)
	fmt.Println(result)
}
