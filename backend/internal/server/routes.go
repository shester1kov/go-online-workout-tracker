package server

import (
	"backend/internal/appmiddlewares"
	"backend/internal/handlers"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

func SetupRoutes(handlers *handlers.Handlers, appmiddlewares *appmiddlewares.AppMiddlewares) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(appmiddlewares.AppCorsMiddleware.Handler())
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("API is running"))
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", handlers.HealthHandler.Check)

		r.Post("/register", handlers.AuthHandler.Register)
		r.Post("/login", handlers.AuthHandler.Login)
		r.Get("/swagger/*", httpSwagger.WrapHandler)

		r.Route("/oauth/fatsecret", func(r chi.Router) {
			r.Get("/callback", handlers.FatSecretAuthHandler.Callback)
		})

		r.Group(func(r chi.Router) {
			r.Use(appmiddlewares.AppAuthMiddlreware.AuthMiddleware())

			r.Post("/logout", handlers.AuthHandler.Logout)

			r.Get("/connect/fatsecret", handlers.FatSecretAuthHandler.ConnectFatSecret)

			r.Route("/nutritions", func(r chi.Router) {
				r.Get("/{date}", handlers.NutritionHandler.GetDailyNutrition)
			})

			r.Route("/exercises", func(r chi.Router) {
				r.Get("/{id}", handlers.ExerciseHandler.GetExercise)
				r.Get("/", handlers.ExerciseHandler.GetExercises)

				r.Group(func(r chi.Router) {
					r.Use(appmiddlewares.AppRoleMiddleware.RoleMiddleware("admin", "moderator"))

					r.Post("/", handlers.ExerciseHandler.CreateExercise)
					r.Put("/{id}", handlers.ExerciseHandler.UpdateExercise)
					r.Delete("/{id}", handlers.ExerciseHandler.DeleteExercise)
				})
			})

			r.Route("/categories", func(r chi.Router) {
				r.Get("/{id}", handlers.CategoryHandler.GetCategory)
				r.Get("/", handlers.CategoryHandler.GetCategories)

				r.Group(func(r chi.Router) {
					r.Use(appmiddlewares.AppRoleMiddleware.RoleMiddleware("admin"))

					r.Post("/", handlers.CategoryHandler.CreateCategory)
					r.Put("/{id}", handlers.CategoryHandler.UpdateCategory)
					r.Delete("/{id}", handlers.CategoryHandler.DeleteCategory)
				})
			})

			r.Route("/users", func(r chi.Router) {
				r.Get("/me", handlers.UserHandler.GetCurrentUser)
				r.Post("/{id}/roles", handlers.UserHandler.AddRoleToUser)
				r.Get("/{id}/roles", handlers.UserHandler.GetUserRoles)
			})

			r.Route("/workouts", func(r chi.Router) {
				r.Get("/{id}/exercises/{workoutExerciseID}", handlers.WorkoutExerciseHandler.GetExerciseByWorkoutID)
				r.Put("/{id}/exercises/{workoutExerciseID}", handlers.WorkoutExerciseHandler.UpdateExerciseInWorkout)
				r.Delete("/{id}/exercises/{workoutExerciseID}", handlers.WorkoutExerciseHandler.DeleteExerciseByWorkoutID)

				r.Post("/{id}/exercises", handlers.WorkoutExerciseHandler.AddExerciseToWorkout)
				r.Get("/{id}/exercises", handlers.WorkoutExerciseHandler.GetExercisesByWorkoutID)

				r.Get("/{id}", handlers.WorkoutHandler.GetWorkoutByUserID)
				r.Put("/{id}", handlers.WorkoutHandler.UpdateWorkoutByUserID)
				r.Delete("/{id}", handlers.WorkoutHandler.DeleteWorkoutByUserID)

				r.Post("/", handlers.WorkoutHandler.CreateWorkout)
				r.Get("/", handlers.WorkoutHandler.GetWorkoutsByUserID)

			})

			r.Route("/foods", func(r chi.Router) {
				r.Get("/{date}", handlers.FoodHandler.GetFood)
				r.Post("/", handlers.FoodHandler.AddFood)
			})
		})
	})

	return r
}
