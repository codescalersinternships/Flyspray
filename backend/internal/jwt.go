package internal

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims is the struct that contains the claims of the JWT
type Claims struct {
	ID string
	jwt.RegisteredClaims
}

// GenerateRefreshToken generates a refresh token
func GenerateToken(secret, userID string, expirationDate time.Time) (string, error) {

	tokenClaims := Claims{
		ID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationDate),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)

	return token.SignedString([]byte(secret))

}

// ValidateToken check that token is valid
func ValidateToken(secret, signedToken string) (*Claims, error) {

	token, err := jwt.ParseWithClaims(signedToken, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
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
