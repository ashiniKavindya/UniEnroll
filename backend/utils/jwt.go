package utils

import (
	"os"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
)

// GenerateToken creates a JWT for the given user ID and role.
func GenerateToken(userID, role string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "devsecret"
	}

	claims := jwt.MapClaims{
		"sub":  userID,
		"role": role,
		"exp":  time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
