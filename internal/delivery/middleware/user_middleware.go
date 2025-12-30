package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/celpung/gocleanarch/pkg/services"
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

// AuthMiddleware:
// - dependency (jwtSvc, secret) DI-INJECT dari luar
// - allowedRoles optional
func AuthMiddleware(jwtSvc *services.JwtService, secret string, allowedRoles ...Role) func(http.Handler) http.Handler {
	// siapkan set role biar O(1) ceknya
	roleSet := map[Role]struct{}{}
	for _, r := range allowedRoles {
		roleSet[r] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 1) Extract token
			tokStr, err := jwtSvc.ExtractBearerToken(r.Header.Get("Authorization"))
			if err != nil {
				writeJSONError(w, http.StatusUnauthorized, "Unauthorized")
				return
			}

			// 2) Parse -> claims
			claims, err := jwtSvc.ParseToken(tokStr, secret)
			if err != nil {
				writeJSONError(w, http.StatusUnauthorized, "Unauthorized")
				return
			}

			// 3) Normalize role
			userRole := Role(strings.ToUpper(strings.TrimSpace(claims.Role)))

			// 4) Role guard (optional)
			if len(roleSet) > 0 {
				if _, ok := roleSet[userRole]; !ok {
					writeJSONError(w, http.StatusForbidden, "Forbidden")
					return
				}
			}

			// 5) Set context (gaya kamu)
			ctx := context.WithValue(r.Context(), ctxKeyID, claims.ID)
			ctx = context.WithValue(ctx, ctxKeyEmail, claims.Email)
			ctx = context.WithValue(ctx, ctxKeyRole, string(userRole))

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Helpers untuk handler/usecase adapter
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
