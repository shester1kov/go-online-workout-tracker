package db

import (
	"backend/internal/config"
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(envs *config.Envs) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", envs.RedisHost, envs.RedisPort),
		Password: envs.RedisPassword,
		DB:       0,
	})

	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return client, nil
}
