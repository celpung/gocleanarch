package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/celpung/gocleanarch/infrastructure/environment"
	"github.com/golang-jwt/jwt/v4"
)

type Role string

const (
	Super Role = "SUPER"
	Admin Role = "ADMIN"
	User  Role = "USER"
)

type contextKey string

const (
	ContextKeyUserID contextKey = "userID"
	ContextKeyEmail  contextKey = "email"
	ContextKeyRole   contextKey = "role"
)

type Claims struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.RegisteredClaims
}

func writeJSONError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]any{"message": msg})
}

func getBearerToken(r *http.Request) (string, error) {
	h := strings.TrimSpace(r.Header.Get("Authorization"))
	if h == "" {
		return "", errors.New("missing Authorization header")
	}
	parts := strings.Fields(h) // e.g. ["Bearer", "<token>"]
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", errors.New("invalid Authorization format")
	}
	return parts[1], nil
}

func AuthMiddleware(requiredRole Role, next http.HandlerFunc) http.HandlerFunc {
	secret := []byte(environment.Env.JWT_SECRET)

	return func(w http.ResponseWriter, r *http.Request) {
		tokStr, err := getBearerToken(r)
		if err != nil {
			writeJSONError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokStr, claims, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return secret, nil
		})
		if err != nil || !token.Valid {
			writeJSONError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		if claims.ExpiresAt != nil && !claims.ExpiresAt.After(time.Now().Add(-30*time.Second)) {
			writeJSONError(w, http.StatusUnauthorized, "Token expired")
			return
		}
		if claims.NotBefore != nil && claims.NotBefore.After(time.Now().Add(30*time.Second)) {
			writeJSONError(w, http.StatusUnauthorized, "Token not valid yet")
			return
		}

		userRole := Role(strings.ToUpper(strings.TrimSpace(claims.Role)))
		if userRole != requiredRole {
			writeJSONError(w, http.StatusForbidden, "Forbidden access: Unauthorized")
			return
		}

		ctx := context.WithValue(r.Context(), ContextKeyUserID, claims.ID)
		ctx = context.WithValue(ctx, ContextKeyEmail, claims.Email)
		ctx = context.WithValue(ctx, ContextKeyRole, string(userRole))

		next(w, r.WithContext(ctx))
	}
}

func UserFromContext(ctx context.Context) (id, email string, role Role, ok bool) {
	idVal, ok1 := ctx.Value(ContextKeyUserID).(string)
	emVal, ok2 := ctx.Value(ContextKeyEmail).(string)
	roVal, ok3 := ctx.Value(ContextKeyRole).(string)
	if !ok1 || !ok2 || !ok3 {
		return "", "", "", false
	}
	return idVal, emVal, Role(roVal), true
}

func UserIDFromContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(ContextKeyUserID).(string)
	return id, ok
}

func UserEmailFromContext(ctx context.Context) (string, bool) {
	email, ok := ctx.Value(ContextKeyEmail).(string)
	return email, ok
}

func UserRoleFromContext(ctx context.Context) (Role, bool) {
	roleStr, ok := ctx.Value(ContextKeyRole).(string)
	return Role(roleStr), ok
}
