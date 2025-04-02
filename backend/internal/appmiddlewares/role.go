package appmiddlewares

import (
	"backend/internal/services"
	"backend/internal/utils"
	"net/http"
)

type AppRoleMiddlreware struct {
	service *services.UserService
}

func NewAppRoleMiddleware(service *services.UserService) *AppRoleMiddlreware {
	return &AppRoleMiddlreware{service: service}
}

func (m *AppRoleMiddlreware) RoleMiddleware(allowedRoleNames ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		allowedSet := make(map[string]struct{})
		for _, allowedRoleName := range allowedRoleNames {
			allowedSet[allowedRoleName] = struct{}{}
		}

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			userID, ok := ctx.Value("user_id").(int)
			if !ok || userID == 0 {
				utils.JSONError(w, "Token or claims not found", http.StatusUnauthorized)
				return
			}

			hasAccess := false

			roles, err := m.service.GetUserRoles(ctx, userID)
			if err != nil {
				utils.JSONError(w, "User have no roles", http.StatusNotFound)
			}

			for _, role := range *roles {
				if _, ok := allowedSet[role.Name]; ok {
					hasAccess = true
					break
				}
			}

			if !hasAccess {
				utils.JSONError(w, "Access denied: insufficient rights", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
