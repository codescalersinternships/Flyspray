package internal

import (
	"os"
	"time"

	"github.com/codescalersinternships/Flyspray/models"
	"github.com/golang-jwt/jwt"
)

// max age in minutes
const tokenMaxAge = 15

func GenerateToken(user models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Minute * tokenMaxAge ).Unix(),
	})

	return token.SignedString([]byte(os.Getenv("SECRET")))
}
