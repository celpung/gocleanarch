package middlewares

import (
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/celpung/gocleanarch/configs/role"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// JWT middleware function with role-based access control
func JWTMiddleware(requiredRole role.Role) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := ctx.GetHeader("Authorization")
		if tokenString == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
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
			return []byte(os.Getenv("JWT_TOKEN")), nil
		})

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Unauthorized",
				"error":   err.Error(),
			})
			return
		}

		// Check if the token is valid
		if !token.Valid {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"success": false,
				"message": "Unauthorized",
				"error":   err,
			})
			return
		}

		userRoleClaim, ok := token.Claims.(jwt.MapClaims)["role"].(float64)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Forbidden access!",
				"error":   "Role claim is not match",
			})
			return
		}

		userRole := role.Role(userRoleClaim)
		if userRole < requiredRole {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Forbidden access!",
				"error":   "Unauthorized",
			})
			return
		}

		// Set the authenticated user in the context
		ctx.Set("userID", token.Claims.(jwt.MapClaims)["id"])
		ctx.Set("email", token.Claims.(jwt.MapClaims)["email"])
		ctx.Set("role", token.Claims.(jwt.MapClaims)["role"])

		// Call the next middleware/handler function in the chain
		ctx.Next()
	}
}
