package app

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/codescalersinternships/Flyspray/internal"
	"github.com/codescalersinternships/Flyspray/models"

	// "github.com/codescalersinternships/Flyspray/models"
	"github.com/gin-gonic/gin"
)

type CustomResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

type CustomError struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
}

func (a *App) Signup(ctx *gin.Context) {
	var user models.User
	err := json.NewDecoder(ctx.Request.Body).Decode(&user)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, CustomError{
			Success: false,
			Error:   err.Error(),
		})
		return
	}
	fmt.Println("request body", user)
	err = user.Validate()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, CustomError{
			Success: false,
			Error:   fmt.Sprintf("validation error: %s",err.Error()),
		})
		return
	}

	hash, err := internal.HashPassword(user.Password)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, CustomError{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	user.Password = hash

	fmt.Println(user)

	// uuidV4 := uuid.New()
	// user.ID = uuidV4.String()

	// user.Verification_code = a.client.GenerateVerificationCode()

	// user, err = a.client.CreateUser(user)
	// if err != nil {
	// 	ctx.JSON(http.StatusInternalServerError, CustomError{
	// 		Success: false,
	// 		Error:   err.Error(),
	// 	})
	// 	return
	// }

	fmt.Println(user)
	// send email
}
