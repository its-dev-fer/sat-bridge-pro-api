package middleware

import (
	"app/src/config"
	"app/src/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// Authenticate is a simplified version of Auth that only checks authentication without role requirements
func Authenticate() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))

		if token == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "Please authenticate")
		}

		userID, err := utils.VerifyToken(token, config.JWTSecret, config.TokenTypeAccess)
		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "Please authenticate")
		}

		// Store userID in context for use in controllers
		c.Locals("userId", userID)

		return c.Next()
	}
}

