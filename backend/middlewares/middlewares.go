package middleware

import (
	"net/http"
	"strings"

	"github.com/codescalersinternships/Flyspray/internal"
	"github.com/gin-gonic/gin"
)

// RequireAuth is a middleware that checks if the user is authenticated
func RequireAuth(ctx *gin.Context, secret string) {
	token := ctx.GetHeader("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")

	if token == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"err": "token not found"})
	}
	claims, err := internal.ValidateToken(secret, token)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"err": "token is invalid"})
		ctx.Abort()
		return
	}

	ctx.Set("user_id", claims.ID)

	ctx.Next()
}
