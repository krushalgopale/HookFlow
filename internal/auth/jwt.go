package auth

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateAccessToken(userID string) (string, error) {
	secret := os.Getenv("JWT_SECRET")

	// Defining Payload/Claims
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
		"iat":     time.Now().Unix(),
	}

	// Defining Algo 
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	// Token/Key Generation using secret code
	return token.SignedString([]byte(secret))
}
