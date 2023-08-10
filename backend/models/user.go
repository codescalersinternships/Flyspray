package models

import (
	"errors"
	"math/rand"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/validator.v2"
)

const verificationCodeLength = 6

// User is the model for the user table
type User struct {
	ID                string `gorm:"primaryKey"`
	Name              string `json:"name"`
	Email             string `json:"email" gorm:"unique;not null" validate:"regexp=^[0-9a-z]+@[0-9a-z]+(\\.[0-9a-z]+)+$"`
	Password          string `json:"password" gorm:"not null"`
	VerificationCode string `gorm:"unique;not null"`
	Verified          bool   `gorm:"default:false"`
}

// Validate validates the user struct
func (u *User) Validate() error {
	if u.Email == "" || u.Name == "" || u.Password == "" {
		return errors.New("missing data name, email, or password")
	}
	return validator.Validate(u)
}

// CreateUser creates a user
func (db *DBClient) CreateUser(user User) (User, error) {
	uuidV4 := uuid.New()
	user.ID = uuidV4.String()
	result := db.Client.Create(&user)
	return user, result.Error
}

// GetUserByID gets a user by id
func (db *DBClient) GetUserByID(id string) (User, error) {
	var user User

	result := db.Client.Select("id", "email", "name", "verified").First(&user)

	return user, result.Error
}

// UpdateUser updates a user
func (db *DBClient) UpdateUser(user User) error {
	if user.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
		if err != nil {
			return err
		}
		result := db.Client.Model(&user).Updates(User{
			Name:     user.Name,
			Email:    user.Email,
			Password: string(hash),
		})
		return result.Error
	}
	result := db.Client.Model(&user).Updates(User{
		Name:  user.Name,
		Email: user.Email,
	})
	return result.Error

}

// GenerateVerificationCode generates a unique verification code
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
