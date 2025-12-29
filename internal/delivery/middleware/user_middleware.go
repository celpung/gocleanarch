package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/celpung/gocleanarch/internal/usecase/port/dependencies"
)

type ctxKey string

const (
	ctxKeyID    ctxKey = "userID"
	ctxKeyEmail ctxKey = "email"
	ctxKeyRole  ctxKey = "role"
)

type Role string

const (
	Super Role = "SUPER"
	Admin Role = "ADMIN"
	User  Role = "USER"
)

func writeJSONError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"message": msg,
	})
}

func getBearerToken(r *http.Request) (string, error) {
	h := strings.TrimSpace(r.Header.Get("Authorization"))
	if h == "" {
		return "", errors.New("missing Authorization header")
	}
	parts := strings.Fields(h)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", errors.New("invalid Authorization format")
	}
	return parts[1], nil
}

func AuthMiddleware(verifier dependencies.TokenVerifier, allowedRoles ...Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if verifier == nil {
				writeJSONError(w, http.StatusInternalServerError, "server misconfigured")
				return
			}
			tokStr, err := getBearerToken(r)
			if err != nil {
				writeJSONError(w, http.StatusUnauthorized, "Unauthorized")
				return
			}

			claims, err := verifier.Verify(tokStr)
			if err != nil {
				writeJSONError(w, http.StatusUnauthorized, "unauthorized")
				return
			}

			if claims.IsExpired(30 * time.Second) {
				writeJSONError(w, http.StatusUnauthorized, "token expired")
				return
			}
			if claims.NotValidYet(30 * time.Second) {
				writeJSONError(w, http.StatusUnauthorized, "token not valid yet")
				return
			}

			userRole := strings.TrimSpace(claims.Role())
			if !roleAllowed(userRole, allowedRoles) {
				writeJSONError(w, http.StatusForbidden, "forbidden")
				return
			}

			ctx := context.WithValue(r.Context(), ctxKeyID, claims.UserID())
			ctx = context.WithValue(ctx, ctxKeyEmail, claims.Email())
			ctx = context.WithValue(ctx, ctxKeyRole, userRole)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Helper untuk dipakai di handler
func UserFromContext(ctx context.Context) (id, email string, role Role, ok bool) {
	idVal, ok1 := ctx.Value(ctxKeyID).(string)
	emVal, ok2 := ctx.Value(ctxKeyEmail).(string)
	roVal, ok3 := ctx.Value(ctxKeyRole).(string)
	if !ok1 || !ok2 || !ok3 {
		return "", "", "", false
	}
	return idVal, emVal, Role(roVal), true
}

func UserIDFromContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(ctxKeyID).(string)
	return id, ok
}

func UserEmailFromContext(ctx context.Context) (string, bool) {
	email, ok := ctx.Value(ctxKeyEmail).(string)
	return email, ok
}

func UserRoleFromContext(ctx context.Context) (Role, bool) {
	role, ok := ctx.Value(ctxKeyRole).(string)
	return Role(role), ok
}

func roleAllowed(userRole string, allowed []Role) bool {
	if len(allowed) == 0 {
		return true
	}
	userRoleNorm := strings.ToLower(userRole)
	for _, r := range allowed {
		if strings.ToLower(string(r)) == userRoleNorm {
			return true
		}
	}
	return false
}
