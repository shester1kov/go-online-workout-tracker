package repository

import "github.com/jmoiron/sqlx"

type Repositories struct {
	ExerciseRepo        *ExerciseRepository
	CategoryRepo        *CategoryRepository
	UserRepo            *UserRepository
	RoleRepo            *RoleRepository
	DBHeathRepo         *DBHeathRepository
	WorkoutRepo         *WorkoutRepository
	WorkoutExerciseRepo *WorkoutExerciseRepository
	FoodRepository      *FoodRepository
}

func InitRepositories(dbConn *sqlx.DB) *Repositories {
	return &Repositories{
		ExerciseRepo:        NewExerciseRepository(dbConn),
		CategoryRepo:        NewCategoryRepository(dbConn),
		UserRepo:            NewUserRepository(dbConn),
		RoleRepo:            NewRoleRepository(dbConn),
		DBHeathRepo:         NewDBHealthRepository(dbConn),
		WorkoutRepo:         NewWorkoutRepository(dbConn),
		WorkoutExerciseRepo: NewWorkoutExerciseRepository(dbConn),
		FoodRepository:      NewFoodRepository(dbConn),
	}
}
