package app

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/codescalersinternships/Flyspray/internal"
	"github.com/codescalersinternships/Flyspray/models"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

type CustomResponse struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

type CustomError struct {
	Error string `json:"error"`
}

func (a *App) Signup(ctx *gin.Context) {
	var user models.User

	err := json.NewDecoder(ctx.Request.Body).Decode(&user)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, CustomError{
			Error: err.Error(),
		})
		return
	}
	err = user.Validate()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, CustomError{
			Error: fmt.Sprintf("validation error: %s", err.Error()),
		})
		return
	}

	hash, err := internal.HashPassword(user.Password)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, CustomError{
			Error: err.Error(),
		})
		return
	}

	user.Password = hash

	user.Verification_code = a.client.GenerateVerificationCode()

	user, err = a.client.CreateUser(user)

	if err != nil {
		if err.Error() == "UNIQUE constraint failed: users.email" {
			ctx.JSON(http.StatusBadRequest, CustomError{
				Error: "email already exists",
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, CustomError{
			Error: err.Error(),
		})
		return
	}
	fmt.Println(user.Verification_code)

	// send email with verification code

	ctx.JSON(http.StatusCreated, CustomResponse{
		Message: "A verification code has been sent to your email.",
	})
}

func (a *App) Verify(ctx *gin.Context) {

	var requestBody struct {
		VerificationCode string `json:"verification_code"`
	}

	err := json.NewDecoder(ctx.Request.Body).Decode(&requestBody)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, CustomError{
			Error: err.Error(),
		})
		return
	}

	result := a.client.Client.Model(&models.User{}).Where("Verification_code = ?", requestBody.VerificationCode).Update("Verified", true)

	if result.RowsAffected != 1 {
		ctx.JSON(http.StatusBadRequest, CustomError{
			Error: "wrong verification code",
		})
		return
	}

	ctx.JSON(http.StatusOK, CustomResponse{
		Message: "your account has been verified successfully",
	})
}

func (a *App) SignIn(ctx *gin.Context) {

	var requestBody models.User
	err := json.NewDecoder(ctx.Request.Body).Decode(&requestBody)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, CustomError{
			Error: err.Error(),
		})
		return
	}

	var user models.User
	result := a.client.Client.First(&user, "Email = ?", requestBody.Email)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, CustomError{
				Error: "invalid email or password",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, CustomError{
			Error: result.Error.Error(),
		})
		return
	}

	if !user.Verified {
		ctx.JSON(http.StatusForbidden, CustomError{
			Error: "Account not verified",
		})
		return
	}

	isPasswordMatches := internal.CheckPasswordHash(requestBody.Password, user.Password)

	if !isPasswordMatches {
		ctx.JSON(http.StatusNotFound, CustomError{
			Error: "invalid email or password",
		})
		return
	}

	// generate token
	token, err := internal.GenerateToken(user)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, CustomError{
			Error: err.Error(),
		})
		return
	}

	user = models.User{
		Name:     user.Name,
		Email:    user.Email,
		Verified: user.Verified,
	}

	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie("Authorization", token, 60*15, "", "", true, true)
	ctx.JSON(http.StatusOK, CustomResponse{
		Message: "logged in successfully",
		Data: struct {
			AccessToken string      `json:"access_token"`
			User        models.User `json:"user"`
		}{
			AccessToken: token,
			User:        user,
		},
	})

}
