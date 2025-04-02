package appmiddlewares

import (
	"backend/internal/auth"
	"backend/internal/utils"
	"context"
	"net/http"
)

type AppAuthMiddlreware struct {
	jwtManager *auth.JWTManager
}

func NewAppAuthMiddleware(jwtManager *auth.JWTManager) *AppAuthMiddlreware {
	return &AppAuthMiddlreware{jwtManager: jwtManager}
}

func (m *AppAuthMiddlreware) AuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("access_token")
			if err != nil {
				utils.JSONError(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			claims, err := m.jwtManager.Verify(cookie.Value)
			if err != nil {
				utils.JSONError(w, "Invalid token", http.StatusForbidden)
				return
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, "user_id", claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
