package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/Bilal-Cplusoft/sun_ready/internal/service"
)

type contextKey string

const (
	UserIDKey    contextKey = "user_id"
	CompanyIDKey contextKey = "company_id"
)

func AuthMiddleware(authService *service.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header required", http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
				return
			}

			token := parts[1]
			claims, err := authService.ValidateToken(token)
			if err != nil {
				http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, CompanyIDKey, claims.CompanyID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserID(ctx context.Context) (int, bool) {
	userID, ok := ctx.Value(UserIDKey).(int)
	return userID, ok
}

func GetCompanyID(ctx context.Context) (int, bool) {
	companyID, ok := ctx.Value(CompanyIDKey).(int)
	return companyID, ok
}
