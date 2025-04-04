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
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("API is running"))
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", handlers.HealthHandler.Check)

		r.Post("/register", handlers.AuthHandler.Register)
		r.Post("/login", handlers.AuthHandler.Login)
		r.Get("/swagger/*", httpSwagger.WrapHandler)

		r.Group(func(r chi.Router) {
			r.Use(appmiddlewares.AppAuthMiddlreware.AuthMiddleware())

			r.Get("/profile", handlers.UserHandler.Profile)
			r.Post("/logout", handlers.AuthHandler.Logout)

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
		})
	})

	return r
}
