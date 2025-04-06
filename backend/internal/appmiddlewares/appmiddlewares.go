package appmiddlewares

import (
	"backend/internal/auth"
	"backend/internal/config"
	"backend/internal/services"
)

type AppMiddlewares struct {
	AppAuthMiddlreware *AppAuthMiddlreware
	AppRoleMiddleware  *AppRoleMiddlreware
	AppCorsMiddleware  *AppCorsMiddleware
}

func InitAppMiddlewares(jwtManager *auth.JWTManager, services *services.Services, envs *config.Envs) *AppMiddlewares {
	return &AppMiddlewares{
		AppAuthMiddlreware: NewAppAuthMiddleware(jwtManager),
		AppRoleMiddleware:  NewAppRoleMiddleware(services.UserService),
		AppCorsMiddleware: NewAppCorsMiddleware(
			[]string{envs.FrontendUrl},
			[]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			[]string{
				"Authorization",
				"Content-Type",
				"Accept",
				"X-Requested-With",
				"X-CSRF-Token",
				"Content-Disposition",
			},
			true,
			300,
			true,
		),
	}
}
