package middleware

import (
	"net/http"
	"strings"

	"github.com/codescalersinternships/Flyspray/internal"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// RequireAuth is a middleware that checks if the user is authenticated
func RequireAuth(secret string) func(*gin.Context) {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")
		token = strings.TrimPrefix(token, "Bearer ")

		if token == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"err": "token not found"})
		}
		claims, err := internal.ValidateToken(secret, token)

		if err != nil {
			log.Error().Err(err).Send()
			ctx.JSON(http.StatusUnauthorized, gin.H{"err": "token is invalid"})
			ctx.Abort()
			return
		}

		ctx.Set("user_id", claims.ID)

		ctx.Next()
	}
}
