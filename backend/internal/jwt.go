package internal

import (
	"errors"
	"os"
	"time"

	"github.com/codescalersinternships/Flyspray/models"
	"github.com/golang-jwt/jwt/v5"
)

// Claims is the struct that contains the claims of the JWT
type Claims struct {
	Email    string
	Name     string
	ID       string
	Verified bool
	jwt.RegisteredClaims
}

// GenerateAccessToken generates an access token
func GenerateAccessToken(user models.User) (string, error) {

	expirationDate := time.Now().Add(time.Minute * 15)
	tokenClaims := Claims{
		Email:    user.Email,
		Name:     user.Name,
		ID:       user.ID,
		Verified: user.Verified,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationDate),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)

	return token.SignedString([]byte(os.Getenv("SECRET")))
}

// GenerateRefreshToken generates a refresh token
func GenerateRefreshToken(user models.User) (string, error) {
	expirationDate := time.Now().Add(time.Hour * 24 * 3)

	tokenClaims := Claims{
		ID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationDate),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)

	return refreshToken.SignedString([]byte(os.Getenv("SECRET")))

}

// ValidateToken check that token is valid
func ValidateToken(signedToken string) (*Claims, error) {

	token, err := jwt.ParseWithClaims(signedToken, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET")), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)

	if !ok {
		return nil, errors.New("invalid token")
	}

	if time.Now().After((*claims.ExpiresAt).Time) {
		return nil, errors.New("token expired")
	}

	return claims, nil
}
