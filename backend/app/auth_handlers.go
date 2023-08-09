package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/codescalersinternships/Flyspray/internal"
	"github.com/codescalersinternships/Flyspray/models"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

func (a *App) Signup(ctx *gin.Context) (interface{}, Response) {
	var user models.User

	err := json.NewDecoder(ctx.Request.Body).Decode(&user)

	if err != nil {
		return nil, BadRequest(err)
	}
	err = user.Validate()
	if err != nil {
		return nil, BadRequest(fmt.Errorf("validation error: %v", err))
	}

	hash, err := internal.HashPassword(user.Password)

	if err != nil {
		return nil, BadRequest(err)
	}

	user.Password = hash

	user.Verification_code = a.client.GenerateVerificationCode()

	user, err = a.client.CreateUser(user)

	if err != nil {
		if err.Error() == "UNIQUE constraint failed: users.email" {
			return nil, Conflict(errors.New("email already exists"))
		}
		return nil, BadRequest(err)
	}

	err = internal.SendEmail(user.Email, user.Verification_code)
	if err != nil {

		return nil, InternalServerError(err)
	}

	message := "A verification code has been sent to your email."

	return ResponseMsg{Message: message}, Created()
}

func (a *App) Verify(ctx *gin.Context) (interface{}, Response) {

	var requestBody struct {
		VerificationCode string `json:"verification_code"`
	}

	err := json.NewDecoder(ctx.Request.Body).Decode(&requestBody)

	if err != nil {
		return nil, BadRequest(err)
	}

	result := a.client.Client.Model(&models.User{}).Where("Verification_code = ?", requestBody.VerificationCode).Update("Verified", true)

	if result.RowsAffected != 1 {
		return nil, BadRequest(errors.New("wrong verification code"))
	}

	msg := "your account has been verified successfully"
	return ResponseMsg{Message: msg}, Ok()

}

func (a *App) SignIn(ctx *gin.Context) (interface{}, Response) {

	var requestBody models.User
	err := json.NewDecoder(ctx.Request.Body).Decode(&requestBody)
	if err != nil {

		return nil, BadRequest(err)
	}

	var user models.User
	result := a.client.Client.First(&user, "Email = ?", requestBody.Email)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, NotFound(errors.New("invalid email or password"))
		}
		return nil, InternalServerError(result.Error)
	}

	if !user.Verified {
		return nil, Forbidden(errors.New("account not verified"))
	}

	isPasswordMatches := internal.CheckPasswordHash(requestBody.Password, user.Password)

	if !isPasswordMatches {

		return nil, NotFound(errors.New("invalid email or password"))
	}

	// generate tokens
	accessToken, err := internal.GenerateAccessToken(user)

	if err != nil {
		return nil, InternalServerError(err)
	}

	refreshToken, err := internal.GenerateRefreshToken(user)

	if err != nil {
		return nil, InternalServerError(err)
	}

	user = models.User{
		Name:     user.Name,
		Email:    user.Email,
		Verified: user.Verified,
	}

	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie("Authorization", accessToken, 60*15, "", "", true, true)

	return ResponseMsg{Message: "logged in successfully", Data: struct {
		AccessToken  string      `json:"access_token"`
		RefreshToken string      `json:"refresh_token"`
		User         models.User `json:"user"`
	}{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}}, Ok()
}

func (a *App) UpdateUser(ctx *gin.Context) (interface{}, Response) {
	userPayload, exists := ctx.Get("user")

	if !exists {
		return nil, UnAuthorized(errors.New("user not found"))
	}

	var user models.User

	if err := json.Unmarshal([]byte(userPayload.(string)), &user); err != nil {
		return nil, InternalServerError(err)
	}

	user, err := a.client.GetUserById(user.ID)
	if err != nil {
		return nil, NotFound(errors.New("user not found"))
	}

	err = json.NewDecoder(ctx.Request.Body).Decode(&user)
	if err != nil {
		return nil, BadRequest(err)
	}

	err = a.client.UpdateUser(user)

	if err != nil {
		return nil, InternalServerError(err)
	}

	return ResponseMsg{Message: "updated successfully"}, Created()
}

func (a *App) GetUser(ctx *gin.Context) (interface{}, Response) {

	userPayload, exists := ctx.Get("user")

	if !exists {
		return nil, NotFound(errors.New("user not found"))
	}

	var user models.User

	if err := json.Unmarshal([]byte(userPayload.(string)), &user); err != nil {
		return nil, InternalServerError(err)
	}

	user, err := a.client.GetUserById(user.ID)

	if err != nil {
		return nil, NotFound(errors.New("user not found"))
	}

	return ResponseMsg{Message: "found", Data: user}, Ok()
}

func (a *App) RefreshToken(ctx *gin.Context) (interface{}, Response) {

	var body struct {
		RefreshToken string `json:"refresh_token"`
	}

	err := json.NewDecoder(ctx.Request.Body).Decode(&body)
	if err != nil {
		return nil, BadRequest(errors.New("token not found"))
	}

	claims, err := internal.ValidateToken(body.RefreshToken)

	if err != nil {
		return nil, UnAuthorized(err)
	}

	user, err := a.client.GetUserById(claims.ID)

	if err != nil {
		return nil, NotFound(errors.New("user not found"))
	}

	accessToken, err := internal.GenerateAccessToken(user)

	if err != nil {
		return nil, InternalServerError(fmt.Errorf("error refreshing token: %v", err))
	}

	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie("Authorization", accessToken, 60*15, "", "", true, true)

	return ResponseMsg{
		Message: "token has been refreshed successfully",
		Data: struct {
			AccessToken string `json:"access_token"`
		}{
			AccessToken: accessToken,
		},
	}, Created()

}
