package middlewares

import (
	"errors"
	"net/http"
	"strings"

	"github.com/celpung/gocleanarch/infrastructure/environment"
	"github.com/celpung/gocleanarch/infrastructure/role"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// JWT middleware function with role-based access control
func AuthMiddleware(requiredRole role.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Token not found!",
				"error":   "Unauthorized",
			})
			return
		}

		tokenString = strings.Replace(tokenString, "bearer ", "", 1)
		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Verify the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid token signing method")
			}

			// Return the secret key used to sign the token
			return []byte(environment.Env.JWT_SECRET), nil
		})

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Unauthorized",
				"error":   err.Error(),
			})
			return
		}

		// Check if the token is valid
		if !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false,
				"message": "Unauthorized",
				"error":   err,
			})
			return
		}

		userRoleClaim, ok := token.Claims.(jwt.MapClaims)["role"].(float64)
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Forbidden access!",
				"error":   "Role claim is not match",
			})
			return
		}

		userRole := role.Role(userRoleClaim)
		if userRole < requiredRole {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Forbidden access!",
				"error":   "Unauthorized",
			})
			return
		}

		// Set the authenticated user in the context
		c.Set("userID", token.Claims.(jwt.MapClaims)["id"])
		c.Set("email", token.Claims.(jwt.MapClaims)["email"])
		c.Set("role", token.Claims.(jwt.MapClaims)["role"])

		// Call the next middleware/handler function in the chain
		c.Next()
	}
}
