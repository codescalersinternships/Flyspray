package app

import (
	"encoding/json"
	"errors"
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

type verifyBody struct {
	VerificationCode int    `json:"verification_code"`
	Email            string `json:"email"`
}

type emailInput struct {
	Email string `json:"email" binding:"required"`
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
		VerificationCodeTimeout: time.Now().Add(time.Second * time.Duration(a.config.MailSender.Timeout)),
	}

	err = user.Validate()
	if err != nil {
		log.Error().Err(err).Send()
		return nil, BadRequest(errors.New("invalid data"))
	}

	hash, err := internal.HashPassword([]byte(user.Password))

	if err != nil {
		log.Error().Err(err).Send()
		return nil, InternalServerError(errInternalServerError)
	}

	user.Password = string(hash)

	user.VerificationCode = a.DB.GenerateVerificationCode()

	user, err = a.DB.CreateUser(user)

	var sqliteErr sqlite3.Error

	if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {

		user, err := a.DB.GetUserByEmail(requestBody.Email)
		if err != nil {
			log.Error().Err(err).Send()
			return nil, InternalServerError(errInternalServerError)
		}
		if !user.Verified {
			verificationCode := a.DB.GenerateVerificationCode()
			err = a.DB.UpdateVerificationCode(user.ID, verificationCode, a.config.MailSender.Timeout)
			if err != nil {
				log.Error().Err(err).Send()
				return nil, InternalServerError(errInternalServerError)
			}

			mailSubject, mailBody := internal.VerifyMailContent(verificationCode)
			err = internal.SendEmail(a.config.MailSender.SendGridKey, a.config.MailSender.Email, user.Email, mailSubject, mailBody)
			if err != nil {
				log.Error().Err(err).Send()
				return nil, InternalServerError(errInternalServerError)
			}
			return ResponseMsg{Message: "your email exists but not verified. check you mailbox for verification code"}, nil
		}
		return nil, Conflict(errors.New("email already exists and verified"))
	}
	if err != nil {
		log.Error().Err(err).Send()

		return nil, InternalServerError(errInternalServerError)
	}

	mailSubject, mailBody := internal.VerifyMailContent(user.VerificationCode)
	err = internal.SendEmail(a.config.MailSender.SendGridKey, a.config.MailSender.Email, user.Email, mailSubject, mailBody)
	if err != nil {
		log.Error().Err(err).Send()

		return nil, InternalServerError(errInternalServerError)
	}

	message := "A verification code has been sent to your email."

	return ResponseMsg{Message: message}, Created()
}

func (a *App) verify(ctx *gin.Context) (interface{}, Response) {

	var requestBody verifyBody

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
		return nil, BadRequest(errors.New("user already verified"))
	}

	if user.VerificationCode != requestBody.VerificationCode {
		return nil, BadRequest(errors.New("wrong verification code"))
	}

	if user.VerificationCodeTimeout.Before(time.Now()) {
		return nil, BadRequest(errors.New("verification code has expired"))
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
	accessToken, err := internal.GenerateToken(a.config.JWT.Secret, user.ID, time.Now().Add(time.Minute*time.Duration(a.config.JWT.Timeout)))

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

	claims, err := internal.ValidateToken(a.config.JWT.Secret, token)

	if err != nil {
		log.Error().Err(err).Send()

		return nil, UnAuthorized(err)
	}

	user, err := a.DB.GetUserByID(claims.ID)

	if err != nil {
		log.Error().Err(err).Send()

		return nil, NotFound(errors.New("user is not found"))
	}

	accessToken, err := internal.GenerateToken(a.config.JWT.Secret, user.ID, time.Now().Add(time.Minute*time.Duration(a.config.JWT.Timeout)))

	if err != nil {
		log.Error().Err(err).Send()

		return nil, InternalServerError(errInternalServerError)
	}

	refreshToken, err := internal.GenerateToken(a.config.JWT.Secret, user.ID, time.Now().Add(time.Hour*time.Duration(a.config.JWT.Timeout)))

	if err != nil {
		log.Error().Err(err).Send()

		return nil, InternalServerError(errInternalServerError)
	}

	return ResponseMsg{
		Message: "token has been refreshed successfully",
		Data: struct {
			AccessToken  string `json:"access_token"`
			RefreshToken string `json:"refresh_token"`
		}{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}, Created()

}

func (a *App) forgetPassword(ctx *gin.Context) (interface{}, Response) {
	var input emailInput

	if err := ctx.BindJSON(&input); err != nil {
		log.Error().Err(err).Send()
		return nil, BadRequest(errors.New("input data is invalid"))
	}

	user, err := a.DB.GetUserByEmail(input.Email)

	if err == gorm.ErrRecordNotFound {
		log.Error().Err(err).Send()
		return nil, NotFound(errors.New("email does not exist"))
	}

	if err != nil {
		log.Error().Err(err).Send()
		return nil, InternalServerError(errInternalServerError)
	}

	if !user.Verified {
		return nil, BadRequest(errors.New("your account is not verified, please verify your account first"))
	}

	verificationCode := a.DB.GenerateVerificationCode()

	if err := a.DB.UpdateVerificationCode(user.ID, verificationCode, a.config.MailSender.Timeout); err != nil {
		log.Error().Err(err).Send()
		return nil, InternalServerError(errInternalServerError)
	}

	mailSubject, mailBody := internal.ForgetPasswordMailContent(verificationCode)
	if err := internal.SendEmail(a.config.MailSender.SendGridKey, a.config.MailSender.Email, user.Email, mailSubject, mailBody); err != nil {
		log.Error().Err(err).Send()
		return nil, InternalServerError(errInternalServerError)
	}

	return ResponseMsg{Message: "forget password code has been sent to your email"}, Ok()
}
