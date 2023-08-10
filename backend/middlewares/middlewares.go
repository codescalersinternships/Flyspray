package middleware

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/codescalersinternships/Flyspray/internal"
	"github.com/codescalersinternships/Flyspray/models"
	"github.com/gin-gonic/gin"
)

// RequireAuth is a middleware that checks if the user is authenticated
func RequireAuth(ctx *gin.Context) {
	tokenString, err := ctx.Cookie("Authorization")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errors.New("token not found"))
		ctx.Abort()
		return
	}

	claims, err := internal.ValidateToken(tokenString)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, err)
		ctx.Abort()
		return
	}
	var user models.User

	user.Email = claims.Email
	user.ID = claims.ID
	user.Name = claims.Name
	user.Verified = claims.Verified

	userEncoded, err := json.Marshal(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		ctx.Abort()
		return
	}
	ctx.Set("user", string(userEncoded))

	ctx.Next()
}
