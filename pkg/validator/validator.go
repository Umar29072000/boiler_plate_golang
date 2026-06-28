package validator

import (
	"regexp"
	"strings"
)

// ValidateEmail validates email format
func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// ValidatePassword validates password strength
func ValidatePassword(password string) bool {
	// Minimum 8 characters
	return len(password) >= 8
}

// ValidateRequired validates required fields
func ValidateRequired(fields map[string]string) []string {
	var errors []string
	for field, value := range fields {
		if strings.TrimSpace(value) == "" {
			errors = append(errors, field+" is required")
		}
	}
	return errors
}
