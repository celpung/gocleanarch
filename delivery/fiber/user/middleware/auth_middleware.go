package middleware

import (
	"errors"
	"strings"
	"time"

	"github.com/celpung/gocleanarch/infrastructure/environment"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type Role string

const (
	Super Role = "SUPER"
	Admin Role = "ADMIN"
	User  Role = "USER"
)

func AuthMiddleware(allowedRoles ...Role) fiber.Handler {
	secret := []byte(environment.Env.JWT_SECRET)

	return func(c *fiber.Ctx) error {
		tokenString, err := getBearerTokenFiber(c)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Unauthorized",
			})
		}

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return secret, nil
		})
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Unauthorized",
			})
		}

		if v, ok := claims["exp"].(float64); ok {
			if time.Now().After(time.Unix(int64(v), 0)) {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"success": false,
					"message": "Token expired",
				})
			}
		}
		if v, ok := claims["nbf"].(float64); ok {
			if time.Now().Before(time.Unix(int64(v), 0)) {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"success": false,
					"message": "Token not valid yet",
				})
			}
		}

		userRole := extractRoleString(claims["role"])
		if len(allowedRoles) > 0 {
			authorized := false
			for _, r := range allowedRoles {
				if strings.EqualFold(userRole, string(r)) {
					authorized = true
					break
				}
			}
			if !authorized {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
					"success": false,
					"message": "Forbidden access!",
				})
			}
		}

		if idStr, ok := claims["id"].(string); ok {
			c.Locals("userID", idStr)
		}
		if emailStr, ok := claims["email"].(string); ok {
			c.Locals("email", emailStr)
		}
		c.Locals("role", strings.ToUpper(userRole))

		return c.Next()
	}
}

func getBearerTokenFiber(c *fiber.Ctx) (string, error) {
	h := strings.TrimSpace(c.Get("Authorization"))
	if h == "" {
		return "", errors.New("missing Authorization header")
	}
	parts := strings.Fields(h) // "Bearer <token>"
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", errors.New("invalid Authorization format")
	}
	return parts[1], nil
}

func extractRoleString(v any) string {
	switch r := v.(type) {
	case string:
		return strings.ToUpper(strings.TrimSpace(r))
	case float64:
		switch int(r) {
		case 1:
			return string(User)
		case 2:
			return string(Admin)
		case 3:
			return string(Super)
		default:
			return ""
		}
	default:
		return ""
	}
}

func UserFromFiberCtx(c *fiber.Ctx) (id, email string, role Role, ok bool) {
	idVal := c.Locals("userID")
	emVal := c.Locals("email")
	roVal := c.Locals("role")

	idStr, ok1 := idVal.(string)
	emStr, ok2 := emVal.(string)
	roStr, ok3 := roVal.(string)
	if !ok1 || !ok2 || !ok3 {
		return "", "", "", false
	}
	return idStr, emStr, Role(roStr), true
}

func UserIDFromFiberCtx(c *fiber.Ctx) (string, bool) {
	if v := c.Locals("userID"); v != nil {
		if s, ok := v.(string); ok {
			return s, true
		}
	}
	return "", false
}

func UserEmailFromFiberCtx(c *fiber.Ctx) (string, bool) {
	if v := c.Locals("email"); v != nil {
		if s, ok := v.(string); ok {
			return s, true
		}
	}
	return "", false
}

func UserRoleFromFiberCtx(c *fiber.Ctx) (Role, bool) {
	if v := c.Locals("role"); v != nil {
		if s, ok := v.(string); ok {
			return Role(s), true
		}
	}
	return "", false
}
