package middleware

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/celpung/gocleanarch/infrastructure/environment"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type Role string

const (
	Super Role = "SUPER"
	Admin Role = "ADMIN"
	User  Role = "USER"
)

func AuthMiddleware(allowedRoles ...Role) gin.HandlerFunc {
	secret := []byte(environment.Env.JWT_SECRET)

	return func(c *gin.Context) {
		tokenString, err := getBearerTokenGin(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
			return
		}

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return secret, nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
			return
		}

		if v, ok := claims["exp"].(float64); ok {
			if time.Now().After(time.Unix(int64(v), 0)) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Token expired"})
				return
			}
		}
		if v, ok := claims["nbf"].(float64); ok {
			if time.Now().Before(time.Unix(int64(v), 0)) {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Token not valid yet"})
				return
			}
		}

		userRole := extractRoleString(claims["role"])
		if len(allowedRoles) > 0 {
			ok := false
			for _, r := range allowedRoles {
				if strings.EqualFold(userRole, string(r)) {
					ok = true
					break
				}
			}
			if !ok {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"success": false, "message": "Forbidden access!"})
				return
			}
		}

		if idStr, ok := claims["id"].(string); ok {
			c.Set("userID", idStr)
		}
		if emailStr, ok := claims["email"].(string); ok {
			c.Set("email", emailStr)
		}
		c.Set("role", strings.ToUpper(userRole))

		// 6) Lanjut
		c.Next()
	}
}

func getBearerTokenGin(c *gin.Context) (string, error) {
	h := strings.TrimSpace(c.GetHeader("Authorization"))
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

func UserFromGinContext(c *gin.Context) (id, email string, role Role, ok bool) {
	idVal, ok1 := c.Get("userID")
	emVal, ok2 := c.Get("email")
	roVal, ok3 := c.Get("role")
	if !ok1 || !ok2 || !ok3 {
		return "", "", "", false
	}
	idStr, _ := idVal.(string)
	emStr, _ := emVal.(string)
	roStr, _ := roVal.(string)
	return idStr, emStr, Role(roStr), true
}

func UserIDFromGinContext(c *gin.Context) (string, bool) {
	if v, ok := c.Get("userID"); ok {
		if s, ok2 := v.(string); ok2 {
			return s, true
		}
	}
	return "", false
}

func UserEmailFromGinContext(c *gin.Context) (string, bool) {
	if v, ok := c.Get("email"); ok {
		if s, ok2 := v.(string); ok2 {
			return s, true
		}
	}
	return "", false
}

func UserRoleFromGinContext(c *gin.Context) (Role, bool) {
	if v, ok := c.Get("role"); ok {
		if s, ok2 := v.(string); ok2 {
			return Role(s), true
		}
	}
	return "", false
}
