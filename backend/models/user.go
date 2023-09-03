package models

import (
	"errors"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"gopkg.in/validator.v2"
	"gorm.io/gorm"
)

// User is the model for the user table
type User struct {
	ID                             string    `json:"id" gorm:"primaryKey"`
	Name                           string    `json:"name"`
	Email                          string    `json:"email" gorm:"unique;not null" validate:"regexp=^[0-9a-z]+@[0-9a-z]+(\\.[0-9a-z]+)+$"`
	Password                       string    `json:"password" gorm:"not null"`
	VerificationCode               int       `json:"verification_code" gorm:"unique"`
	Verified                       bool      `json:"verified" gorm:"default:false"`
	VerificationCodeExpirationTime time.Time `json:"verification_code_expiration_time"`
}

// Validate validates the user struct
func (u *User) Validate() error {
	if u.Email == "" || u.Name == "" || u.Password == "" {
		return errors.New("missing data name, email, or password")
	}
	return validator.Validate(u)
}

// BeforeCreate generates a new uuid
func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	uuidV4 := uuid.New()
	if err != nil {
		return err
	}

	user.ID = uuidV4.String()
	return
}

// CreateUser creates a user
func (db *DBClient) CreateUser(user User) (User, error) {
	result := db.Client.Create(&user)
	return user, result.Error
}

// GetUserByID gets a user by id
func (db *DBClient) GetUserByID(id string) (User, error) {
	user := User{}
	return user, db.Client.First(&user, "id = ?", id).Error
}

// UpdateUser updates a user
func (db *DBClient) UpdateUser(user User) error {

	result := db.Client.Model(&user).Updates(User{
		Name:     user.Name,
		Password: user.Password,
	})
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return result.Error
}

// GenerateVerificationCode generates a unique verification code
func (db *DBClient) GenerateVerificationCode() int {

	var verificationCode int

	for {
		verificationCode = rand.Intn(900000) + 100000
		var count int64
		db.Client.Model(&User{}).Where("Verification_code = ?", verificationCode).Count(&count)

		if count == 0 {
			return verificationCode
		}
	}

}

func (db *DBClient) UpdateVerificationCode(userID string, newVerificationCode int, timeout int) error {
	return db.Client.Model(&User{}).Where("id = ?", userID).Updates(User{
		VerificationCode:               newVerificationCode,
		VerificationCodeExpirationTime: time.Now().Add(time.Second * time.Duration(timeout)),
	}).Error
}

func (db *DBClient) VerifyUser(userID string) error {
	return db.Client.Model(&User{}).Where("id = ?", userID).Update("verified", true).Error
}

func (db *DBClient) GetUserByEmail(email string) (User, error) {
	var user User
	result := db.Client.First(&user, "Email = ?", email)

	return user, result.Error

}
