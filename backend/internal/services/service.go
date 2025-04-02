package services

import (
	"backend/internal/auth"
	"backend/internal/repository"

	"github.com/redis/go-redis/v9"
)

type Services struct {
	ExerciseService *ExerciseService
	CategoryService *CategoryService
	UserService     *UserService
	AuthService     *AuthService
	HealthService   *HealthService
}

func InitServices(repos *repository.Repositories, redis *redis.Client, jwtManager *auth.JWTManager) *Services {
	return &Services{
		ExerciseService: NewExerciseService(repos.ExerciseRepo, repos.CategoryRepo, redis),
		CategoryService: NewCategoryService(repos.CategoryRepo, redis),
		UserService:     NewUserService(repos.UserRepo, repos.RoleRepo),
		AuthService:     NewAuthService(repos.UserRepo, jwtManager),
		HealthService:   NewHealthService(repos.DBHeathRepository, redis),
	}
}
