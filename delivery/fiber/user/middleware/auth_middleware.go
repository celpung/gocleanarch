package middleware

import (
	"errors"
	"strings"

	"github.com/celpung/gocleanarch/infrastructure/environment"
	"github.com/celpung/gocleanarch/infrastructure/role"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

// AuthMiddleware is a JWT middleware with role-based access control
func AuthMiddleware(requiredRole role.Role) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get token from header
		tokenString := c.Get("Authorization")
		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Token not found!",
				"error":   "Unauthorized",
			})
		}

		// Strip Bearer prefix
		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
		tokenString = strings.Replace(tokenString, "bearer ", "", 1)

		// Parse JWT
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid token signing method")
			}
			return []byte(environment.Env.JWT_SECRET), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Unauthorized",
				"error":   err.Error(),
			})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"message": "Invalid token claims",
			})
		}

		// Get role from claims
		roleClaim, ok := claims["role"].(float64)
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"message": "Role claim is missing or invalid",
			})
		}

		userRole := role.Role(roleClaim)
		if userRole < requiredRole {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"message": "Forbidden access!",
				"error":   "Unauthorized",
			})
		}

		// Set context values
		c.Locals("userID", claims["id"])
		c.Locals("email", claims["email"])
		c.Locals("role", claims["role"])

		// Continue to the next handler
		return c.Next()
	}
}
