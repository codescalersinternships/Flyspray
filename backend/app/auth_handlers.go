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
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

type signupBody struct {
	Name            string `json:"name"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
	Verified        bool   `json:"verified"`
}

type signinBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (a *App) signup(ctx *gin.Context) (interface{}, Response) {

	var requestBody signupBody

	err := json.NewDecoder(ctx.Request.Body).Decode(&requestBody)

	if err != nil {
		log.Error().Err(err).Msg("")
		return nil, BadRequest(errors.New("input data is invalid"))
	}

	if requestBody.Password != requestBody.ConfirmPassword {
		return nil, BadRequest(errors.New("passwords do not match"))
	}

	user := models.User{
		Name:                    requestBody.Name,
		Email:                   requestBody.Email,
		Password:                requestBody.Password,
		Verified:                requestBody.Verified,
		VerificationCodeTimeout: time.Now().Add(time.Hour * 2),
	}

	err = user.Validate()
	if err != nil {
		log.Error().Err(err).Msg("")
		return nil, BadRequest(fmt.Errorf("validation error: %w", err))
	}

	hash, err := internal.HashPassword(user.Password)

	if err != nil {
		log.Error().Err(err).Msg("")
		return nil, BadRequest(fmt.Errorf("failed to hash password: %w", err))
	}

	user.Password = hash

	user.VerificationCode = a.client.GenerateVerificationCode()
	fmt.Println(user.VerificationCode)

	user, err = a.client.CreateUser(user)

	if err != nil {
		log.Error().Err(err).Msg("")
		var sqliteErr sqlite3.Error

		if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {

			user, err := a.client.GetUserByEmail(requestBody.Email)
			if err != nil {
				log.Error().Err(err).Msg("")
				return nil, InternalServerError(err)
			}
			if !user.Verified {
				return ResponseMsg{Message: "your email exists but not verified. check you mailbox for verification code"}, Conflict(errors.New("email already exists"))
			}
			return nil, Conflict(errors.New("email already exists"))
		}
		return nil, BadRequest(err)
	}

	err = internal.SendEmail(user.Email, user.VerificationCode)
	if err != nil {
		log.Error().Err(err).Msg("")

		return nil, InternalServerError(fmt.Errorf("failed to send verification mail: %w", err))
	}

	message := "A verification code has been sent to your email."

	return ResponseMsg{Message: message}, Created()
}

func (a *App) verify(ctx *gin.Context) (interface{}, Response) {

	var requestBody struct {
		VerificationCode int `json:"verification_code"`
	}

	err := json.NewDecoder(ctx.Request.Body).Decode(&requestBody)

	if err != nil {
		log.Error().Err(err).Msg("")

		return nil, BadRequest(errors.New("input data is invalid"))
	}

	msg, err := a.client.VerifyUser(requestBody.VerificationCode)

	if err != nil {
		log.Error().Err(err).Msg("")
		return nil, BadRequest(err)
	}

	return ResponseMsg{Message: msg}, Ok()

}

func (a *App) signIn(ctx *gin.Context) (interface{}, Response) {

	var requestBody signinBody
	err := json.NewDecoder(ctx.Request.Body).Decode(&requestBody)
	if err != nil {
		log.Error().Err(err).Msg("")

		return nil, BadRequest(errors.New("input data is invalid"))
	}

	user, err := a.client.GetUserByEmail(requestBody.Email)

	if err != nil {
		log.Error().Err(err).Msg("")

		if err == gorm.ErrRecordNotFound {
			return nil, NotFound(errors.New("wrong email or password"))
		}
		return nil, BadRequest(err)
	}

	isPasswordMatches := internal.CheckPasswordHash(requestBody.Password, user.Password)

	if !isPasswordMatches {

		return nil, NotFound(errors.New("wrong email or password"))
	}

	if !user.Verified {
		return nil, Forbidden(errors.New("account is not verified yet, please check the verification email in your inbox"))
	}
	// generate tokens
	accessToken, err := internal.GenerateToken(user.ID, time.Now().Add(time.Minute*15))

	if err != nil {
		log.Error().Err(err).Msg("")

		return nil, InternalServerError(fmt.Errorf("failed to generate token:%w ", err))
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

	_, err := a.client.GetUserByID(userID.(string))
	if err != nil {
		log.Error().Err(err).Msg("")

		return nil, NotFound(errors.New("user is not found"))
	}

	var requestBody struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	err = json.NewDecoder(ctx.Request.Body).Decode(&requestBody)
	if err != nil {
		log.Error().Err(err).Msg("")

		return nil, BadRequest(err)
	}

	user := models.User{
		Name:  requestBody.Name,
		Email: requestBody.Email,
		ID:    userID.(string),
	}

	err = a.client.UpdateUser(user)

	if err != nil {
		log.Error().Err(err).Msg("")

		return nil, InternalServerError(err)
	}

	return ResponseMsg{Message: "updated successfully"}, Ok()
}

func (a *App) getUser(ctx *gin.Context) (interface{}, Response) {

	userID, exists := ctx.Get("user_id")

	if !exists {
		return nil, NotFound(errors.New("user is not found"))
	}

	user, err := a.client.GetUserByID(userID.(string))

	if err != nil {
		log.Error().Err(err).Msg("")

		return nil, NotFound(errors.New("user is not found"))
	}

	return ResponseMsg{Message: "found", Data: user}, Ok()
}

func (a *App) refreshToken(ctx *gin.Context) (interface{}, Response) {

	token := ctx.GetHeader("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")

	claims, err := internal.ValidateToken(token)

	if err != nil {
		log.Error().Err(err).Msg("")

		return nil, UnAuthorized(err)
	}

	user, err := a.client.GetUserByID(claims.ID)

	if err != nil {
		log.Error().Err(err).Msg("")

		return nil, NotFound(errors.New("user is not found"))
	}

	accessToken, err := internal.GenerateToken(user.ID, time.Now().Add(time.Minute*15))

	if err != nil {
		log.Error().Err(err).Msg("")

		return nil, InternalServerError(fmt.Errorf("error refreshing token: %v", err))
	}

	refreshToken, err := internal.GenerateToken(user.ID, time.Now().Add(time.Hour*24*3))

	if err != nil {
		log.Error().Err(err).Msg("")

		return nil, InternalServerError(fmt.Errorf("error refreshing token: %v", err))
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
