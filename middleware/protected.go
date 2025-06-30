package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/ridoystarlord/bondhu-pay/utils"
)

// Protected returns middleware that verifies JWT and attaches user ID to context
func Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Extract Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return utils.Unauthorized(c,"Missing Authorization header")
		}

		// Must be Bearer <token>
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			return utils.Unauthorized(c,"Invalid Authorization header format")
		}
		tokenString := parts[1]

		// Validate and parse token
		userID, err := utils.ParseJWT(tokenString)
		if err != nil {
			return utils.Unauthorized(c,"Invalid or expired token")
		}

		// Store userID in locals for later handlers
		c.Locals("userID", userID)

		// Continue to the next handler
		return c.Next()
	}
}
