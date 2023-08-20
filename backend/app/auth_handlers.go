package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/codescalersinternships/Flyspray/internal"
	"github.com/codescalersinternships/Flyspray/models"
	"github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
	"gopkg.in/validator.v2"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

const timeout = 120
const tokenTimeout = 15
const apiKey = ""
const apiEmail = ""
const secret = ""

type signupBody struct {
	Name            string `json:"name"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

type signinBody struct {
	Email    string `json:"email" validate:"nonzero"`
	Password string `json:"password" validate:"nonzero"`
}

func (a *App) signup(ctx *gin.Context) (interface{}, Response) {

	var requestBody signupBody

	err := json.NewDecoder(ctx.Request.Body).Decode(&requestBody)

	if err != nil {
		log.Error().Err(err).Send()
		return nil, BadRequest(errors.New("input data is invalid"))
	}

	if requestBody.Password != requestBody.ConfirmPassword {
		return nil, BadRequest(errors.New("passwords do not match"))
	}

	user := models.User{
		Name:                    requestBody.Name,
		Email:                   requestBody.Email,
		Password:                requestBody.Password,
		VerificationCodeTimeout: time.Now().Add(time.Second * time.Duration(timeout)),
	}

	err = user.Validate()
	if err != nil {
		log.Error().Err(err).Send()
		return nil, BadRequest(errors.New("invalid data"))
	}

	hash, err := internal.HashPassword([]byte(user.Password))

	if err != nil {
		log.Error().Err(err).Send()
		return nil, BadRequest(errors.New("failed to hash password"))
	}

	user.Password = string(hash)

	user.VerificationCode = a.DB.GenerateVerificationCode()
	fmt.Println(user.VerificationCode)

	user, err = a.DB.CreateUser(user)

	var sqliteErr sqlite3.Error

	if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {

		user, err := a.DB.GetUserByEmail(requestBody.Email)
		if err != nil {
			log.Error().Err(err).Send()
			return nil, InternalServerError(errInternalServerError)
		}
		if !user.Verified {
			verifivationCode := a.DB.GenerateVerificationCode()
			err = a.DB.UpdateVerificationCode(user.ID, verifivationCode, timeout)
			if err != nil {
				log.Error().Err(err).Send()
				return "", InternalServerError(errInternalServerError)
			}

			err = internal.SendEmail(apiKey, apiEmail, user.Email, verifivationCode)
			if err != nil {
				log.Error().Err(err).Send()
				return "", InternalServerError(errInternalServerError)
			}
			return ResponseMsg{Message: "your email exists but not verified. check you mailbox for verification code"}, nil
		}
		return nil, Conflict(errors.New("email already exists and verified"))
	}
	if err != nil {
		log.Error().Err(err).Send()

		return nil, InternalServerError(errInternalServerError)
	}

	err = internal.SendEmail(apiKey, apiEmail, user.Email, user.VerificationCode)
	if err != nil {
		log.Error().Err(err).Send()

		return nil, InternalServerError(errInternalServerError)
	}

	message := "A verification code has been sent to your email."

	return ResponseMsg{Message: message}, Created()
}

func (a *App) verify(ctx *gin.Context) (interface{}, Response) {

	var requestBody struct {
		VerificationCode int    `json:"verification_code"`
		Email            string `json:"email"`
	}

	err := json.NewDecoder(ctx.Request.Body).Decode(&requestBody)

	if err != nil {
		log.Error().Err(err).Send()

		return nil, BadRequest(errors.New("input data is invalid"))
	}

	user, err := a.DB.GetUserByEmail(requestBody.Email)
	if err == gorm.ErrRecordNotFound {
		return nil, NotFound(errors.New("email does not exist"))
	}
	if err != nil {
		log.Error().Err(err).Send()
		return nil, InternalServerError(errInternalServerError)
	}

	if user.Verified {
		return "", BadRequest(errors.New("user already verified"))
	}

	if user.VerificationCode != requestBody.VerificationCode {
		return "", BadRequest(errors.New("wrong verification code"))
	}

	fmt.Println(user.VerificationCodeTimeout)
	fmt.Println(time.Now().Add(time.Duration(0)))

	if user.VerificationCodeTimeout.Before(time.Now()) {
		return "", BadRequest(errors.New("verification code has expired"))
	}

	err = a.DB.VerifyUser(user.ID)

	if err != nil {
		log.Error().Err(err).Send()
		return nil, BadRequest(errInternalServerError)
	}

	return ResponseMsg{Message: "verified"}, Ok()

}

func (a *App) signIn(ctx *gin.Context) (interface{}, Response) {

	var requestBody signinBody
	err := json.NewDecoder(ctx.Request.Body).Decode(&requestBody)
	if err != nil {
		log.Error().Err(err).Send()

		return nil, BadRequest(errors.New("input data is invalid"))
	}

	err = validator.Validate(requestBody)
	if err != nil {
		log.Error().Err(err).Send()
		return nil, BadRequest(errors.New("invalid input data"))
	}
	user, err := a.DB.GetUserByEmail(requestBody.Email)

	if err == gorm.ErrRecordNotFound {
		return nil, NotFound(errors.New("wrong email or password"))
	}

	if err != nil {
		log.Error().Err(err).Send()
		return nil, InternalServerError(errInternalServerError)
	}

	isPasswordMatches := internal.CheckPasswordHash([]byte(user.Password), requestBody.Password)

	if !isPasswordMatches {

		return nil, NotFound(errors.New("wrong email or password"))
	}

	if !user.Verified {
		return nil, Forbidden(errors.New("account is not verified yet, please check the verification email in your inbox"))
	}
	// generate tokens
	accessToken, err := internal.GenerateToken(secret, user.ID, time.Now().Add(time.Minute*time.Duration(tokenTimeout)))

	if err != nil {
		log.Error().Err(err).Send()

		return nil, InternalServerError(errInternalServerError)
	}

	return ResponseMsg{Message: "logged in successfully", Data: struct {
		AccessToken string `json:"access_token"`
	}{
		AccessToken: accessToken,
	}}, Ok()
}

func (a *App) updateUser(ctx *gin.Context) (interface{}, Response) {
	userID, exists := ctx.Get("user_id")

	if !exists {
		return nil, UnAuthorized(errors.New("user is not found"))
	}

	var requestBody struct {
		Name string `json:"name"`
	}
	err := json.NewDecoder(ctx.Request.Body).Decode(&requestBody)
	if err != nil {
		log.Error().Err(err).Send()

		return nil, BadRequest(errors.New("invalid input data"))
	}

	user := models.User{
		Name: requestBody.Name,
		ID:   userID.(string),
	}

	err = a.DB.UpdateUser(user)

	if err == gorm.ErrRecordNotFound {
		return nil, BadRequest(errors.New("user does not exist"))
	}
	if err != nil {
		log.Error().Err(err).Send()

		return nil, InternalServerError(errInternalServerError)
	}

	return ResponseMsg{Message: "user has been updated successfully"}, Ok()
}

func (a *App) getUser(ctx *gin.Context) (interface{}, Response) {

	userID, exists := ctx.Get("user_id")

	if !exists {
		return nil, NotFound(errors.New("user is not found"))
	}

	user, err := a.DB.GetUserByID(userID.(string))

	if err == gorm.ErrRecordNotFound {
		return nil, NotFound(errors.New("user is not found"))
	}
	if err != nil {
		log.Error().Err(err).Send()
		return nil, InternalServerError(errInternalServerError)
	}

	return ResponseMsg{Message: "user is found", Data: user}, Ok()
}

func (a *App) refreshToken(ctx *gin.Context) (interface{}, Response) {

	token := ctx.GetHeader("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")
	if token == "" {
		return nil, UnAuthorized(errors.New("token is required"))
	}

	claims, err := internal.ValidateToken(secret, token)

	if err != nil {
		log.Error().Err(err).Send()

		return nil, UnAuthorized(err)
	}

	user, err := a.DB.GetUserByID(claims.ID)

	if err != nil {
		log.Error().Err(err).Send()

		return nil, NotFound(errors.New("user is not found"))
	}

	accessToken, err := internal.GenerateToken(secret, user.ID, time.Now().Add(time.Minute*time.Duration(tokenTimeout)))

	if err != nil {
		log.Error().Err(err).Send()

		return nil, InternalServerError(errInternalServerError)
	}

	refreshToken, err := internal.GenerateToken(secret, user.ID, time.Now().Add(time.Hour*time.Duration(tokenTimeout)))

	if err != nil {
		log.Error().Err(err).Send()

		return nil, InternalServerError(errInternalServerError)
	}

	return ResponseMsg{
		Message: "token has been refreshed successfully",
		Data: struct {
			AccessToken  string `json:"access_token"`
			RefreshToken string `json:"resfresh_token"`
		}{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}, Created()

}
