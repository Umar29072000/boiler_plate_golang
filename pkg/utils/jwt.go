package utils

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret string

// JWTClaims represents JWT token claims
type JWTClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// InitJWT initializes JWT secret
func InitJWT(secret string) {
	jwtSecret = secret
}

// ValidateToken validates JWT token and returns claims
func ValidateToken(tokenString string) (*JWTClaims, error) {
	if jwtSecret == "" {
		return nil, errors.New("JWT secret not initialized")
	}

	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
