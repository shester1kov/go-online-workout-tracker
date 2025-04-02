package handlers

import "backend/internal/services"

type Handlers struct {
	ExerciseHandler *ExerciseHandler
	CategoryHandler *CategoryHandler
	UserHandler     *UserHandler
	AuthHandler     *AuthHandler
	HealthHandler   *HealthHandler
}

func InitHandlers(services *services.Services) *Handlers {
	return &Handlers{
		ExerciseHandler: NewExerciseHandler(services.ExerciseService),
		CategoryHandler: NewCategoryHandler(services.CategoryService),
		UserHandler:     NewUserHandler(services.UserService),
		AuthHandler:     NewAuthHandler(services.AuthService),
		HealthHandler:   NewHealthHandler(services.HealthService),
	}
}
