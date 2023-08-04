package models

import (
	"math/rand"

	"gopkg.in/validator.v2"
)

const verificationCodeLength = 6

type User struct {
	ID                string `gorm:"primary_key;<-:false"`
	Name              string `json:"name" validate:"required"`
	Email             string `json:"email" gorm:"unique;not null" validate:"required regexp=^[0-9a-z]+@[0-9a-z]+(\\.[0-9a-z]+)+$"`
	Password          string `json:"password" gorm:"not null" validate:"required"`
	Verification_code string `gorm:"unique;not null"`
	Verified          bool   `gorm:"default:false"`
}

func (u *User) Validate() error {
	return validator.Validate(u)
}

func (db *DBClient) CreateUser(user User) (User, error) {

	result := db.Client.Create(&user)
	return user, result.Error
}

func (db *DBClient) GetUser(id string) error {
	var user User
	result := db.Client.First(&user, id)

	return result.Error
}

func (db *DBClient) GenerateVerificationCode() string {

	charsSet := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

	verificationCode := make([]byte, verificationCodeLength)

	for {
		for i := range verificationCode {
			verificationCode[i] = charsSet[rand.Intn(len(charsSet)-1)]
		}

		var count int64
		db.Client.Model(&User{}).Where("verificationCode = ?", string(verificationCode)).Count(&count)

		if count == 0 {
			return string(verificationCode)
		}
	}

}
