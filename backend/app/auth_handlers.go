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
	VerificationCode int    `json:"verification_code" binding:"required"`
	Email            string `json:"email" binding:"required"`
}

type forgetPasswordBody struct {
	Email string `json:"email" binding:"required"`
}

type changePasswordBody struct {
	Password        string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
}

type updateUserBody struct {
	Name string `json:"name"`
}

// signup creates a new user account and sends a verification code to the user's email
// @Summary Create a new user account
// @Description Creates a new user account and sends a verification code to the user's email
// @Tags Users
// @Accept json
// @Produce json
// @Param request body signupBody true "Signup request body"
// @Success 201 {object} ResponseMsg ""A verification code has been sent to your email"
// @Failure 400 {object} Response "Bad request"
// @Failure 409 {object} Response "Email already exists and verified"
// @Failure 500 {object} Response "Internal server error"
// @Router /user/signup [post]
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
		Name:                           requestBody.Name,
		Email:                          requestBody.Email,
		Password:                       requestBody.Password,
		VerificationCodeExpirationTime: time.Now().Add(time.Second * time.Duration(a.config.MailSender.Timeout)),
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

// verify verifies a user's account using the provided verification code and email
// @Summary Verify a user's account
// @Description Verifies a user's account using the provided verification code and email
// @Tags Users
// @Accept json
// @Produce json
// @Param request body verifyBody true "Verification request body"
// @Success 200 {object} ResponseMsg "verified"
// @Failure 400 {object} Response "Bad request"
// @Failure 404 {object} Response "Email does not exist"
// @Failure 500 {object} Response "Internal server error"
// @Router /user/signup/verify [post]
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

	if user.VerificationCodeExpirationTime.Before(time.Now()) {
		return nil, BadRequest(errors.New("verification code has expired"))
	}

	err = a.DB.VerifyUser(user.ID)

	if err != nil {
		log.Error().Err(err).Send()
		return nil, BadRequest(errInternalServerError)
	}

	return ResponseMsg{Message: "verified"}, Ok()

}

// signIn authenticates a user's credentials and generates access tokens for the user
// @Summary Authenticate a user and generate access tokens
// @Description Authenticates a user's credentials and generates access tokens
// @Tags Users
// @Accept json
// @Produce json
// @Param request body signinBody true "Signin request body"
// @Success 200 {object} ResponseMsg "logged in successfully (AccessToken details is given in a struct in the 'Data' field)"
// @Failure 400 {object} Response "Bad request"
// @Failure 401 {object} Response "Unauthorized"
// @Failure 403 {object} Response "Forbidden"
// @Failure 404 {object} Response "Wrong email or password"
// @Failure 500 {object} Response "Internal server error"
// @Router /user/signin [post]
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

// updateUser updates the user's information
// @Summary Update user information
// @Description Updates the user's information
// @Tags Users
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body updateUserBody true "Update user request body"
// @Success 200 {object} ResponseMsg "user has been updated successfully"
// @Failure 400 {object} Response "Bad request"
// @Failure 401 {object} Response "Unauthorized"
// @Failure 404 {object} Response "User does not exist"
// @Failure 500 {object} Response "Internal server error"
// @Router /user [put]
func (a *App) updateUser(ctx *gin.Context) (interface{}, Response) {

	var requestBody updateUserBody

	userID, exists := ctx.Get("user_id")

	if !exists {
		return nil, UnAuthorized(errors.New("user is not found"))
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

// getUser retrieves the user's information
// @Summary Get user information
// @Description Retrieves the user's information
// @Tags Users
// @Produce json
// @Security Bearer
// @Success 200 {object} ResponseMsg "user is found (User details in the 'Data' field)"
// @Failure 401 {object} Response "Unauthorized"
// @Failure 404 {object} Response "User does not exist"
// @Failure 500 {object} Response "Internal server error"
// @Router /user [get]
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

// refreshToken generates a new access token using the provided refresh token
// @Summary Generate new access token
// @Description Generates a new access token using the provided refresh token
// @Tags Users
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Security Bearer
// @Success 201 {object} ResponseMsg  "token has been refreshed successfully (AccessToken & RefreshToken details are given in a struct in the 'Data' field)"
// @Failure 400 {object} Response "Bad request"
// @Failure 401 {object} Response "Unauthorized"
// @Failure 500 {object} Response "Internal server error"
// @Router /user/refresh-token [post]
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

// forgetPassword requests forget password code to be sent to user's email
// @Summary Request forget password code
// @Description requests forget password code to be sent to user's email
// @Tags Users
// @Accept json
// @Produce json
// @Param request body forgetPasswordBody true "request forget password code request body"
// @Success 200 {object} ResponseMsg  "forget password code has been sent to your email"
// @Failure 400 {object} Response "Bad request"
// @Failure 404 {object} Response "NotFound"
// @Failure 500 {object} Response "Internal server error"
// @Router /user/forget_password [post]
func (a *App) forgetPassword(ctx *gin.Context) (interface{}, Response) {
	var input forgetPasswordBody

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

// verifyForgetPassword verify forget password code then send token
// @Summary verify forget password code
// @Description verify forget password code then send token
// @Tags Users
// @Accept json
// @Produce json
// @Param request body verifyBody true "verify forget password code request body"
// @Success 200 {object} ResponseMsg "verified (AccessToken details is given in a struct in the 'Data' field)"
// @Failure 400 {object} Response "Bad request"
// @Failure 404 {object} Response "NotFound"
// @Failure 500 {object} Response "Internal server error"
// @Router /user/forget_password/verify [post]
func (a *App) verifyForgetPassword(ctx *gin.Context) (interface{}, Response) {
	var input verifyBody

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

	if user.VerificationCode != input.VerificationCode {
		return nil, BadRequest(errors.New("wrong verification code"))
	}

	if user.VerificationCodeExpirationTime.Before(time.Now()) {
		return nil, BadRequest(errors.New("verification code has expired"))
	}

	// generate tokens
	accessToken, err := internal.GenerateToken(a.config.JWT.Secret, user.ID, time.Now().Add(time.Minute*time.Duration(a.config.JWT.Timeout)))

	if err != nil {
		log.Error().Err(err).Send()
		return nil, InternalServerError(errInternalServerError)
	}

	return ResponseMsg{Message: "verified", Data: struct {
		AccessToken string `json:"access_token"`
	}{
		AccessToken: accessToken,
	}}, Ok()
}

// changePassword changes password
// @Summary changes password
// @Description changes password
// @Tags Users
// @Accept json
// @Produce json
// @Param request body changePasswordBody true "change password request body"
// @Param Authorization header string true "Bearer token"
// @Security Bearer
// @Success 200 {object} ResponseMsg "password has been updated successfully"
// @Failure 400 {object} Response "Bad request"
// @Failure 401 {object} Response "UnAuthorized"
// @Failure 500 {object} Response "Internal server error"
// @Router /user/password [put]
func (a *App) changePassword(ctx *gin.Context) (interface{}, Response) {
	var input changePasswordBody

	if err := ctx.BindJSON(&input); err != nil {
		log.Error().Err(err).Send()
		return nil, BadRequest(errors.New("input data is invalid"))
	}

	userID, exists := ctx.Get("user_id")
	if !exists {
		return nil, UnAuthorized(errors.New("authentication is required"))
	}

	if input.Password != input.ConfirmPassword {
		return nil, BadRequest(errors.New("passwords do not match"))
	}

	hashedPassword, err := internal.HashPassword([]byte(input.Password))
	if err != nil {
		log.Error().Err(err).Send()
		return nil, InternalServerError(errInternalServerError)
	}

	user := models.User{ID: userID.(string), Password: string(hashedPassword)}

	if err := a.DB.UpdateUser(user); err != nil {
		log.Error().Err(err).Send()
		return nil, InternalServerError(errInternalServerError)
	}

	return ResponseMsg{Message: "password has been updated successfully"}, Ok()
}
