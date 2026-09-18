package middleware

import (
	"context"
	"net/http"
	"strings"
	"uuid"

	"github.com/jwhite9387/taskflow-api/internal/auth"
)

type contextKey string

const userIDKey contextKey = "userID"

func RequireAuth(jwtManager *auth.JWTManager, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")

		if authorization == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		parts := strings.Fields(authorization)

		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		userID, err := jwtManager.ValidateToken(parts[1])
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		r = r.WithContext(WithUserID(r.Context(), userID))

		next.ServeHTTP(w, r)
	})
}

func UserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	value, ok := ctx.Value(userIDKey).(uuid.UUID)
	if !ok {
		return uuid.Nil(), false
	}
	return value, true
}

func WithUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}
