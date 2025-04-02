package services

import (
	"backend/internal/models"
	"backend/internal/repository"
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type HealthService struct {
	dbRepo *repository.DBHeathRepository
	redis  *redis.Client
}

func NewHealthService(dbRepo *repository.DBHeathRepository, redis *redis.Client) *HealthService {
	return &HealthService{
		dbRepo: dbRepo,
		redis:  redis,
	}
}

func (s *HealthService) Check(ctx context.Context) *models.HealthStatus {
	status := models.HealthStatus{
		Status:    "up",
		Timestamp: time.Now().UTC(),
		Details:   make(map[string]string),
	}

	if err := s.dbRepo.Check(ctx); err != nil {
		status.Status = "down"
		status.Details["database"] = err.Error()
	}

	if _, err := s.redis.Ping(ctx).Result(); err != nil {
		status.Status = "down"
		status.Details["redis"] = err.Error()
	}

	return &status
}
