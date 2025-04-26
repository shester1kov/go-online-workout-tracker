package services

import (
	"backend/internal/auth"
	"backend/internal/clients"
	"backend/internal/oauth"
	"backend/internal/repository"

	"github.com/redis/go-redis/v9"
)

type Services struct {
	ExerciseService        *ExerciseService
	CategoryService        *CategoryService
	UserService            *UserService
	AuthService            *AuthService
	HealthService          *HealthService
	WorkoutSerivce         *WorkoutSerivce
	WorkoutExerciseSerivce *WorkoutExerciseSerivce
	FoodService            *FoodService
	NutritionService       *NutritionService
}

func InitServices(repos *repository.Repositories, redis *redis.Client, jwtManager *auth.JWTManager, clients *clients.Clients, oauth *oauth.Oauth) *Services {
	return &Services{
		ExerciseService:        NewExerciseService(repos.ExerciseRepo, repos.CategoryRepo, redis),
		CategoryService:        NewCategoryService(repos.CategoryRepo, redis),
		UserService:            NewUserService(repos.UserRepo, repos.RoleRepo),
		AuthService:            NewAuthService(repos.UserRepo, jwtManager),
		HealthService:          NewHealthService(repos.DBHeathRepo, redis),
		WorkoutSerivce:         NewWorkoutService(repos.WorkoutRepo),
		WorkoutExerciseSerivce: NewWorkoutExerciseService(repos.WorkoutRepo, repos.WorkoutExerciseRepo, repos.ExerciseRepo),
		FoodService:            NewFoodService(clients.NutritionixClient, repos.FoodRepository),
		NutritionService:       NewNutritionService(repos.FatSecretAuthRepository, oauth.FatSecretAuthClient),
	}
}
