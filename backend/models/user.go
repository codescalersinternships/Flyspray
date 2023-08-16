package models

import (
	"errors"
	"math/rand"
	"time"

	"github.com/codescalersinternships/Flyspray/internal"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gopkg.in/validator.v2"
)

// User is the model for the user table
type User struct {
	ID                      string `gorm:"primaryKey"`
	Name                    string `json:"name"`
	Email                   string `json:"email" gorm:"unique;not null" validate:"regexp=^[0-9a-z]+@[0-9a-z]+(\\.[0-9a-z]+)+$"`
	Password                string `json:"password" gorm:"not null"`
	VerificationCode        int    `gorm:"unique"`
	Verified                bool   `gorm:"default:false"`
	VerificationCodeTimeout time.Time
}

// Validate validates the user struct
func (u *User) Validate() error {
	if u.Email == "" || u.Name == "" || u.Password == "" {
		return errors.New("missing data name, email, or password")
	}
	return validator.Validate(u)
}

// BeforeCreate generates a uuid for the user
func (db *DBClient) BeforeCreate(user User) {
	uuidV4 := uuid.New()
	user.ID = uuidV4.String()
}

// CreateUser creates a user
func (db *DBClient) CreateUser(user User) (User, error) {
	db.BeforeCreate(user)
	result := db.Client.Create(&user)
	return user, result.Error
}

// GetUserByID gets a user by id
func (db *DBClient) GetUserByID(id string) (User, error) {
	var user User

	result := db.Client.First(&user)

	return user, result.Error
}

// UpdateUser updates a user
func (db *DBClient) UpdateUser(user User) error {

	result := db.Client.Model(&user).Updates(User{
		Name:  user.Name,
		Email: user.Email,
	})
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

func (db *DBClient) VerifyUser(verificationCode int) (string, error) {

	var user User
	result := db.Client.First(&user, "verification_code = ?", verificationCode)

	if result.Error != nil {
		return "", result.Error
	}

	if user.ID == "" {
		return "", errors.New("wrong verification code")
	}

	if user.Verified {
		return "", errors.New("user already verified")
	}

	result = db.Client.Model(&User{}).Where("verification_code = ? AND verification_code_Timeout > ?", verificationCode, time.Now()).Update("Verified", true)

	if result.Error != nil {
		log.Error().Err(result.Error).Msg("")
		return "", errors.New("verification code has expired and failed to generate new one")
	}
	if result.RowsAffected != 1 {
		verificationCode := db.GenerateVerificationCode()
		result = db.Client.Model(&User{}).Where("verification_code = ?", verificationCode).Updates(User{VerificationCode: verificationCode, VerificationCodeTimeout: time.Now().Add(time.Hour * 2)})

		if result.Error != nil || result.RowsAffected != 1 {
			log.Error().Err(result.Error).Msg("")
			return "", errors.New("verification code has expired and failed to generate new one")
		}

		err := internal.SendEmail(user.Email, verificationCode)

		if err != nil {
			log.Error().Err(err).Msg("")
			return "", err
		}

		return "your verification code has expired. A new one has been sent to your email", nil

	}
	return "updated succsssfully", nil

}

func (db *DBClient) GetUserByEmail(email string) (User, error) {
	var user User
	result := db.Client.First(&user, "Email = ?", email)

	return user, result.Error

}
