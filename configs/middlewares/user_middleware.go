package middlewares

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/celpung/gocleanarch/configs/role"
	"github.com/golang-jwt/jwt"
)

// JWTMiddleware function with role-based access control for net/http
func JWTMiddleware(requiredRole role.Role, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Token not found!", http.StatusUnauthorized)
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
			http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
			return
		}

		// Check if the token is valid
		if !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userRoleClaim, ok := token.Claims.(jwt.MapClaims)["role"].(float64)
		if !ok {
			http.Error(w, "Forbidden access: Role claim does not match", http.StatusForbidden)
			return
		}

		userRole := role.Role(userRoleClaim)
		if userRole < requiredRole {
			http.Error(w, "Forbidden access: Unauthorized", http.StatusForbidden)
			return
		}

		// Set the authenticated user in the request context
		ctx := r.Context()
		ctx = context.WithValue(ctx, "userID", token.Claims.(jwt.MapClaims)["id"])
		ctx = context.WithValue(ctx, "email", token.Claims.(jwt.MapClaims)["email"])
		ctx = context.WithValue(ctx, "role", token.Claims.(jwt.MapClaims)["role"])

		// Call the next middleware/handler function in the chain
		next(w, r.WithContext(ctx))
	}
}
