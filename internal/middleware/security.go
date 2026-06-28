package middleware

import (
	"github.com/gofiber/fiber/v2"
)

// SecurityHeaders adds security headers to responses (Helmet.js equivalent)
func SecurityHeaders() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// X-XSS-Protection: Enables XSS protection in older browsers
		c.Set("X-XSS-Protection", "1; mode=block")

		// X-Content-Type-Options: Prevents MIME type sniffing
		c.Set("X-Content-Type-Options", "nosniff")

		// X-Frame-Options: Prevents clickjacking
		c.Set("X-Frame-Options", "DENY")

		// X-DNS-Prefetch-Control: Controls DNS prefetching
		c.Set("X-DNS-Prefetch-Control", "off")

		// Referrer-Policy: Controls referrer information
		c.Set("Referrer-Policy", "strict-origin-when-cross-origin")

		// Permissions-Policy: Controls browser features
		c.Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		// Strict-Transport-Security: Enforces HTTPS (only in production)
		// Note: Only enable this if you're using HTTPS
		// c.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

		return c.Next()
	}
}
