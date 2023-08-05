package models

import (
	"math/rand"

	"github.com/google/uuid"
	"gopkg.in/validator.v2"
)

const verificationCodeLength = 6

type User struct {
	ID                string `gorm:"primaryKey;<-:false"`
	Name              string `json:"name" validate:"required"`
	Email             string `json:"email" gorm:"unique;not null"`
	Password          string `json:"password" gorm:"not null"`
	Verification_code string `gorm:"unique;not null"`
	Verified          bool   `gorm:"default:false"`
}

func (u *User) Validate() error {
	return validator.Validate(u)
}

func (db *DBClient) CreateUser(user User) (User, error) {
	uuidV4 := uuid.New()
	user.ID = uuidV4.String()
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
		db.Client.Model(&User{}).Where("Verification_code = ?", string(verificationCode)).Count(&count)

		if count == 0 {
			return string(verificationCode)
		}
	}

}
