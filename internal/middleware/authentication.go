// Package middleware provides reusable middleware components for the application.
// These include authentication, logging, and other request/response processing utilities.
package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

// JWTValidator is an interface that defines a method for validating JWT tokens.
type JWTValidator interface {
	ValidateJWT(token string) (string, error)
}

// RequireAuth is a middleware that enforces authentication for protected routes.
// It validates the JWT token from the request and sets the user ID in the context.
//
// Parameters:
//   - jwt: An implementation of the JWTValidator interface for token validation.
//
// Returns:
//   - fiber.Handler: The middleware handler function.
func RequireAuth(jwt JWTValidator) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Cookies("access_token")
		if token == "" {
			authHeader := c.Get("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				token = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}

		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "authentication required",
			})
		}

		userID, err := jwt.ValidateJWT(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid or expired token",
			})
		}

		// Inject user ID into context
		c.Locals("user_id", userID)

		// Inject into forwarded headers
		c.Request().Header.Set("X-User-ID", userID)

		return c.Next()
	}
}
