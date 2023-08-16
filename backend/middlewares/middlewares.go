package middleware

import (
	"net/http"
	"strings"

	"github.com/codescalersinternships/Flyspray/internal"
	"github.com/codescalersinternships/Flyspray/models"
	"github.com/gin-gonic/gin"
)

// RequireAuth is a middleware that checks if the user is authenticated
func RequireAuth(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")
	token = strings.TrimPrefix(token, "Bearer ")

	claims, err := internal.ValidateToken(token)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"err": err})
		ctx.Abort()
		return
	}
	var user models.User

	user.ID = claims.ID

	ctx.Set("user_id", claims.ID)

	ctx.Next()
}
