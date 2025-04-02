package repository

import (
	"database/sql"
)

type Repositories struct {
	ExerciseRepo      *ExerciseRepository
	CategoryRepo      *CategoryRepository
	UserRepo          *UserRepository
	RoleRepo          *RoleRepository
	DBHeathRepository *DBHeathRepository
}

func InitRepositories(dbConn *sql.DB) *Repositories {
	return &Repositories{
		ExerciseRepo:      NewExerciseRepository(dbConn),
		CategoryRepo:      NewCategoryRepository(dbConn),
		UserRepo:          NewUserRepository(dbConn),
		RoleRepo:          NewRoleRepository(dbConn),
		DBHeathRepository: NewDBHealthRepository(dbConn),
	}
}
