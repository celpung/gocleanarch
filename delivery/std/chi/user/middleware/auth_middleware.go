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

// Typed claims supaya tidak perlu casting-casting MapClaims.
type Claims struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.RegisteredClaims
}

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
	// Pecah by space: "Bearer <token>"
	parts := strings.Fields(h)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", errors.New("invalid Authorization format")
	}
	return parts[1], nil
}

func AuthMiddleware(allowedRoles ...Role) func(http.Handler) http.Handler {
	secret := []byte(environment.Env.JWT_SECRET)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokStr, err := getBearerToken(r)
			if err != nil {
				writeJSONError(w, http.StatusUnauthorized, "Unauthorized")
				return
			}

			claims := &Claims{}
			token, err := jwt.ParseWithClaims(tokStr, claims, func(t *jwt.Token) (interface{}, error) {
				// Pastikan pakai HMAC (HS256/HS384/HS512)
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("unexpected signing method")
				}
				return secret, nil
			})
			if err != nil || !token.Valid {
				writeJSONError(w, http.StatusUnauthorized, "Unauthorized")
				return
			}

			// Validasi waktu (exp/nbf) dengan leeway kecil (opsional).
			if claims.ExpiresAt != nil && !claims.ExpiresAt.After(time.Now().Add(-30*time.Second)) {
				writeJSONError(w, http.StatusUnauthorized, "Token expired")
				return
			}
			if claims.NotBefore != nil && claims.NotBefore.After(time.Now().Add(30*time.Second)) {
				writeJSONError(w, http.StatusUnauthorized, "Token not valid yet")
				return
			}

			// Cek role (normalize uppercase)
			userRole := Role(strings.ToUpper(strings.TrimSpace(claims.Role)))
			if len(allowedRoles) > 0 {
				authorized := false
				for _, r := range allowedRoles {
					if userRole == r {
						authorized = true
						break
					}
				}
				if !authorized {
					writeJSONError(w, http.StatusForbidden, "Forbidden")
					return
				}
			}

			ctx := context.WithValue(r.Context(), ctxKeyID, claims.ID)
			ctx = context.WithValue(ctx, ctxKeyEmail, claims.Email)
			ctx = context.WithValue(ctx, ctxKeyRole, string(userRole))

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
