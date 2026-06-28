package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

// generateRandomToken generates a random token of specified length (internal helper)
func generateRandomToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}

// GenerateVerificationToken generates email verification token (32 bytes = 64 hex chars)
func GenerateVerificationToken() (string, error) {
	return generateRandomToken(32)
}

// GenerateResetToken generates password reset token (32 bytes = 64 hex chars)
func GenerateResetToken() (string, error) {
	return generateRandomToken(32)
}
