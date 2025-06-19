package middlewares

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/celpung/gocleanarch/configs/environment"
	"github.com/celpung/gocleanarch/configs/role"
	"github.com/golang-jwt/jwt/v4"
)

type contextKey string

const (
	ContextKeyUserID contextKey = "userID"
	ContextKeyEmail  contextKey = "email"
	ContextKeyRole   contextKey = "role"
)

// AuthMiddleware verifies JWT token and enforces role-based access control.
func AuthMiddleware(requiredRole role.Role, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Token not found!", http.StatusUnauthorized)
			return
		}

		// Strip Bearer prefix
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		tokenString = strings.TrimPrefix(tokenString, "bearer ")

		// Parse and validate token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Ensure signing method is HMAC
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid token signing method")
			}
			return []byte(environment.Env.JWT_SECRET), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
			return
		}

		// Extract role and verify access
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusForbidden)
			return
		}

		roleFloat, ok := claims["role"].(float64)
		if !ok {
			http.Error(w, "Invalid or missing role claim", http.StatusForbidden)
			return
		}
		userRole := role.Role(roleFloat)
		if userRole < requiredRole {
			http.Error(w, "Forbidden access: Unauthorized", http.StatusForbidden)
			return
		}

		// Inject claims into context using safe keys
		ctx := context.WithValue(r.Context(), ContextKeyUserID, claims["id"])
		ctx = context.WithValue(ctx, ContextKeyEmail, claims["email"])
		ctx = context.WithValue(ctx, ContextKeyRole, claims["role"])

		next(w, r.WithContext(ctx))
	}
}
