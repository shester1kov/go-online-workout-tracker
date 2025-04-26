package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Envs struct {
	Port                    string
	PostgresHost            string
	PostgresPort            string
	PostgresUser            string
	PostgresPassword        string
	PostgresDb              string
	PostgresSslmode         string
	RedisHost               string
	RedisPort               string
	RedisPassword           string
	JWTSecureKey            string
	FrontendUrl             string
	NutritionixAppID        string
	NutritionixAppKey       string
	FatsecretConsumerKey    string
	FatsecretConsumerSecret string
	FatsecretCallbackURL    string
}

func LoadEnvs() *Envs {
	env := os.Getenv("ENV")
	if env != "docker" {
		err := godotenv.Load("../.env")
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	return &Envs{
		Port:                    os.Getenv("PORT"),
		PostgresHost:            os.Getenv("POSTGRES_HOST"),
		PostgresPort:            os.Getenv("POSTGRES_PORT"),
		PostgresUser:            os.Getenv("POSTGRES_USER"),
		PostgresPassword:        os.Getenv("POSTGRES_PASSWORD"),
		PostgresDb:              os.Getenv("POSTGRES_DB"),
		PostgresSslmode:         os.Getenv("POSTGRES_SSLMODE"),
		RedisHost:               os.Getenv("REDIS_HOST"),
		RedisPort:               os.Getenv("REDIS_PORT"),
		RedisPassword:           os.Getenv("REDIS_PASSWORD"),
		JWTSecureKey:            os.Getenv("JWT_SECURE_KEY"),
		FrontendUrl:             os.Getenv("FRONTEND_URL"),
		NutritionixAppID:        os.Getenv("NUTRITIONIX_APP_ID"),
		NutritionixAppKey:       os.Getenv("NUTRITIONIX_APP_KEY"),
		FatsecretConsumerKey:    os.Getenv("FATSECRET_CONSUMER_KEY"),
		FatsecretConsumerSecret: os.Getenv("FATSECRET_CONSUMER_SECRET"),
		FatsecretCallbackURL:    os.Getenv("FATSECRET_CALLBACK_URL"),
	}
}
