package appmiddlewares

import (
	"backend/internal/auth"
	"backend/internal/services"
)

type AppMiddlewares struct {
	AppAuthMiddlreware *AppAuthMiddlreware
	AppRoleMiddleware  *AppRoleMiddlreware
}

func InitAppMiddlewares(jwtManager *auth.JWTManager, services *services.Services) *AppMiddlewares {
	return &AppMiddlewares{
		AppAuthMiddlreware: NewAppAuthMiddleware(jwtManager),
		AppRoleMiddleware:  NewAppRoleMiddleware(services.UserService),
	}
}
