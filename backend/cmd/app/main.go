package main

import (
	"backend/internal/appmiddlewares"
	"backend/internal/auth"
	"backend/internal/clients"
	"backend/internal/config"
	"backend/internal/db"
	"backend/internal/handlers"
	"backend/internal/oauth"
	"backend/internal/repository"
	"backend/internal/server"
	"backend/internal/services"
	"log"

	_ "backend/docs"
)

// @title Online Workout Tracker API
// @version 1.0
// @description This is a API for managing online workout tracker
// @BasePath /api/v1
// @contact.name API Support
// @contact.url https://www.example.com/support
// @contact.email support@example.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
func main() {
	envs, err := config.LoadEnvs("../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	dbConn, err := db.NewConnection(envs)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbConn.Close()

	if err := db.RunMigrations(dbConn, envs); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	redisClient, err := db.NewRedisClient(envs)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisClient.Close()

	repos := repository.InitRepositories(dbConn)
	jwtManager := auth.InitJWTManager(envs)
	clients := clients.InitClients(envs)
	oauth := oauth.InitOauth(envs)
	service := services.InitServices(repos, redisClient, jwtManager, clients, oauth)
	handler := handlers.InitHandlers(service, envs)
	appmiddleware := appmiddlewares.InitAppMiddlewares(jwtManager, service, envs)

	router := server.SetupRoutes(handler, appmiddleware)

	server.StartServer(router, envs.Port)

}
