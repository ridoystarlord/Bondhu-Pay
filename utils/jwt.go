package utils

import (
	"errors"
	"os"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var (
	jwtSecret []byte
	once      sync.Once
)

func loadSecret() {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		panic("JWT_SECRET is not set in environment variables")
	}
	jwtSecret = []byte(secret)
}

// GenerateJWT creates a signed token with user ID
func GenerateJWT(userID string) (string, error) {
	once.Do(loadSecret)

	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(72 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtSecret)
}

// ParseJWT validates and extracts user ID
func ParseJWT(tokenString string) (string, error) {
	once.Do(loadSecret)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return jwtSecret, nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", errors.New("invalid token claims")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", errors.New("user_id not found in token")
	}

	return userID, nil
}
