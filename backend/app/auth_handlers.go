package app

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/codescalersinternships/Flyspray/internal"
	"github.com/codescalersinternships/Flyspray/models"

	// "github.com/google/uuid"

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
	user := models.User{
		ID:                "123", // Assign a valid value to ID
		Name:              "diaa",    // You can leave this empty if you want to bypass the validation
		Email:             "test@test.com",
		Password:          "password123",
		Verification_code: "abcd123",
		Verified:          false,
	}

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
			Error:   fmt.Sprintf("validation error: %s", err.Error()),
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

	user.Verification_code = a.client.GenerateVerificationCode()

	user, err = a.client.CreateUser(user)
	if err != nil {
		if err.Error() == "UNIQUE constraint failed: users.email" {
			ctx.JSON(http.StatusBadRequest, CustomError{
				Success: false,
				Error:   "email already exists",
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, CustomError{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	fmt.Println(user.ID)

}
