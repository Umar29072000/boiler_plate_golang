package middleware

import (
	"html"
	"regexp"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// XSSProtection sanitizes user input to prevent XSS attacks
func XSSProtection() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get request body if exists
		if len(c.Body()) > 0 {
			// Sanitize body
			sanitized := sanitizeInput(string(c.Body()))
			c.Request().SetBody([]byte(sanitized))
		}

		// Sanitize query parameters
		c.Request().URI().QueryArgs().VisitAll(func(key, value []byte) {
			sanitized := sanitizeInput(string(value))
			c.Request().URI().QueryArgs().Set(string(key), sanitized)
		})

		return c.Next()
	}
}

// InputSanitizer provides more aggressive input sanitization
func InputSanitizer() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Sanitize common injection patterns
		if len(c.Body()) > 0 {
			body := string(c.Body())
			
			// Remove SQL injection patterns (defense in depth, GORM handles this)
			body = removeSQLInjectionPatterns(body)
			
			// Remove script tags and event handlers
			body = removeScriptTags(body)
			
			c.Request().SetBody([]byte(body))
		}

		return c.Next()
	}
}

// sanitizeInput escapes HTML and removes dangerous patterns
func sanitizeInput(input string) string {
	// HTML escape
	input = html.EscapeString(input)
	
	// Remove null bytes
	input = strings.ReplaceAll(input, "\x00", "")
	
	return input
}

// removeScriptTags removes script tags and event handlers
func removeScriptTags(input string) string {
	// Remove <script> tags (case insensitive)
	scriptRegex := regexp.MustCompile(`(?i)<script[^>]*>.*?</script>`)
	input = scriptRegex.ReplaceAllString(input, "")
	
	// Remove event handlers (onclick, onerror, etc)
	eventRegex := regexp.MustCompile(`(?i)\s*on\w+\s*=\s*["'][^"']*["']`)
	input = eventRegex.ReplaceAllString(input, "")
	
	// Remove javascript: protocol
	jsRegex := regexp.MustCompile(`(?i)javascript:`)
	input = jsRegex.ReplaceAllString(input, "")
	
	return input
}

// removeSQLInjectionPatterns removes common SQL injection patterns
func removeSQLInjectionPatterns(input string) string {
	// Note: GORM uses parameterized queries which prevent SQL injection
	// This is defense in depth
	
	// Remove common SQL keywords in suspicious contexts
	patterns := []string{
		`(?i)(\s|^)(union|select|insert|update|delete|drop|create|alter|exec|execute)(\s|$)`,
		`(?i)--`,           // SQL comments
		`(?i)/\*.*?\*/`,    // SQL comments
		`(?i);\s*(drop|delete|update|insert)`, // Dangerous statements after semicolon
	}
	
	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		// Replace with empty string or sanitized version
		input = re.ReplaceAllString(input, " ")
	}
	
	return input
}
