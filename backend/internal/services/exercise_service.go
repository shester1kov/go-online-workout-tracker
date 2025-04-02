package services

import (
	"backend/internal/apperrors"
	"backend/internal/models"
	"backend/internal/repository"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

const exerciseCacheKey = "exercises"

type ExerciseService struct {
	exerciseRepo *repository.ExerciseRepository
	categoryRepo *repository.CategoryRepository
	redis        *redis.Client
}

func NewExerciseService(exerciseRepo *repository.ExerciseRepository, categoryRepo *repository.CategoryRepository, redis *redis.Client) *ExerciseService {
	return &ExerciseService{
		exerciseRepo: exerciseRepo,
		categoryRepo: categoryRepo,
		redis:        redis,
	}
}

func (s *ExerciseService) CreateExercise(ctx context.Context, req *models.ExerciseRequest) (*models.Exercise, error) {

	if category, err := s.categoryRepo.GetCategory(ctx, req.CategoryID); category == nil || err != nil {
		return nil, &apperrors.AppError{
			Code:    http.StatusBadRequest,
			Message: "Incorrect category id",
		}
	}

	exercise := &models.Exercise{
		Name:        req.Name,
		Description: req.Description,
		CategoryID:  req.CategoryID,
	}

	err := s.exerciseRepo.CreateExercise(ctx, exercise)
	if err != nil {
		var pgErr *pq.Error
		switch {
		case errors.Is(err, context.Canceled):
			log.Println("Request cancelled:", err)
			return nil, &apperrors.AppError{
				Code:    http.StatusBadRequest,
				Message: "Request cancelled",
			}

		case errors.Is(err, context.DeadlineExceeded):
			log.Println("Deadline exceeded:", err)
			return nil, &apperrors.AppError{
				Code:    http.StatusGatewayTimeout,
				Message: "Request timeout",
			}

		case errors.As(err, &pgErr) && pgErr.Code == apperrors.PgErrUniqueViolation:
			log.Println("Unique violation:", pgErr)
			return nil, &apperrors.AppError{
				Code:    http.StatusConflict,
				Message: "Exercise already exists",
			}

		case errors.As(err, &pgErr) && pgErr.Code == apperrors.PgErrForeignKeyViolation:
			log.Println("Foreign key violation:", pgErr)
			return nil, &apperrors.AppError{
				Code:    http.StatusBadRequest,
				Message: "Incorrect category id",
			}

		default:
			log.Println("Unhandled error:", err)
			return nil, &apperrors.AppError{
				Code:    http.StatusInternalServerError,
				Message: "Failed to create exercise",
			}
		}
	}

	s.redis.Del(ctx, exerciseCacheKey)
	return exercise, nil
}

func (s *ExerciseService) GetExercises(ctx context.Context) (*[]models.Exercise, error) {
	val, err := s.redis.Get(ctx, exerciseCacheKey).Result()
	if err == redis.Nil {
		log.Println("cache not found, getting data from DB")
	} else if err != nil {
		log.Println("error getting data from Redis:", err)
	} else {
		var exercises []models.Exercise

		if err := json.Unmarshal([]byte(val), &exercises); err == nil {
			log.Println("data retrieved from cache")
			return &exercises, nil
		}

		log.Println("error deserializing data from cache:", err)
	}

	exercises, err := s.exerciseRepo.GetExercises(ctx)
	if err != nil {
		switch {
		case errors.Is(err, context.Canceled):
			log.Println("Request cancelled:", err)
			return nil, &apperrors.AppError{
				Code:    http.StatusBadRequest,
				Message: "Request cancelled",
			}
		case errors.Is(err, context.DeadlineExceeded):
			log.Println("Deadline exceeded:", err)
			return nil, &apperrors.AppError{
				Code:    http.StatusGatewayTimeout,
				Message: "Request timeout",
			}
		default:
			log.Println("Unhandled error:", err)
			return nil, &apperrors.AppError{
				Code:    http.StatusInternalServerError,
				Message: "Failed to get exercises",
			}
		}
	}

	if exercises == nil || len(*exercises) == 0 {
		log.Println(err)
		return nil, &apperrors.AppError{
			Code:    http.StatusNotFound,
			Message: "Exercises not found",
		}
	}

	data, err := json.Marshal(exercises)
	if err != nil {
		log.Println("cache serialization error", err)
	} else {
		s.redis.Set(ctx, exerciseCacheKey, data, 10*time.Minute)
		log.Println("data written to cache")
	}

	return exercises, nil
}

func (s *ExerciseService) GetExercise(ctx context.Context, id int) (*models.Exercise, error) {

	exercise, err := s.exerciseRepo.GetExercise(ctx, id)
	if err != nil || exercise == nil {
		switch {
		case errors.Is(err, context.Canceled):
			log.Println("Request cancelled:", err)
			return nil, &apperrors.AppError{
				Code:    http.StatusBadRequest,
				Message: "Request cancelled",
			}

		case errors.Is(err, context.DeadlineExceeded):
			log.Println("Deadline exceeded:", err)
			return nil, &apperrors.AppError{
				Code:    http.StatusGatewayTimeout,
				Message: "Request timeout",
			}

		case errors.Is(err, sql.ErrNoRows):
			log.Println("Category not found:", err)
			return nil, &apperrors.AppError{
				Code:    http.StatusNotFound,
				Message: "Exercise not found",
			}

		default:
			log.Println("Unhandled error:", err)
			return nil, &apperrors.AppError{
				Code:    http.StatusInternalServerError,
				Message: "Failed to get exercise",
			}
		}
	}

	return exercise, nil
}

func (s *ExerciseService) UpdateExercise(ctx context.Context, id int, req *models.ExerciseRequest) (*models.Exercise, error) {

	if category, err := s.categoryRepo.GetCategory(ctx, req.CategoryID); category == nil || err != nil {
		return nil, &apperrors.AppError{
			Code:    http.StatusBadRequest,
			Message: "Incorrect category id",
		}
	}

	exercise := &models.Exercise{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		CategoryID:  req.CategoryID,
	}

	err := s.exerciseRepo.UpdateExercise(ctx, exercise)
	if err != nil {
		var pgErr *pq.Error
		switch {
		case errors.Is(err, context.Canceled):
			log.Println("Request cancelled:", err)
			return nil, &apperrors.AppError{
				Code:    http.StatusBadRequest,
				Message: "Request cancelled",
			}

		case errors.Is(err, context.DeadlineExceeded):
			log.Println("Deadline exceeded:", err)
			return nil, &apperrors.AppError{
				Code:    http.StatusGatewayTimeout,
				Message: "Request timeout",
			}

		case errors.Is(err, sql.ErrNoRows):
			log.Println("Category not found:", err)
			return nil, &apperrors.AppError{
				Code:    http.StatusNotFound,
				Message: "Exercise not found",
			}

		case errors.As(err, &pgErr) && pgErr.Code == apperrors.PgErrUniqueViolation:
			log.Println("Unique violation:", pgErr)
			return nil, &apperrors.AppError{
				Code:    http.StatusConflict,
				Message: "Exercise already exists",
			}

		case errors.As(err, &pgErr) && pgErr.Code == apperrors.PgErrForeignKeyViolation:
			log.Println("Foreign key violation:", pgErr)
			return nil, &apperrors.AppError{
				Code:    http.StatusBadRequest,
				Message: "Incorrect category id",
			}

		default:
			log.Println("Unhandled error:", err)
			return nil, &apperrors.AppError{
				Code:    http.StatusInternalServerError,
				Message: "Failed to update exercise",
			}
		}
	}

	s.redis.Del(ctx, exerciseCacheKey)
	return exercise, nil
}

func (s *ExerciseService) DeleteExercise(ctx context.Context, id int) error {
	rowsAffected, err := s.exerciseRepo.DeleteExercise(ctx, id)

	if err != nil {
		switch {
		case errors.Is(err, context.Canceled):
			log.Println("Request cancelled:", err)
			return &apperrors.AppError{
				Code:    http.StatusBadRequest,
				Message: "Request cancelled",
			}
		case errors.Is(err, context.DeadlineExceeded):
			log.Println("Deadline exceeded:", err)
			return &apperrors.AppError{
				Code:    http.StatusGatewayTimeout,
				Message: "Request timeout",
			}
		default:
			log.Println("Unhandled error:", err)
			return &apperrors.AppError{
				Code:    http.StatusInternalServerError,
				Message: "Failed to get categories",
			}
		}
	}

	if rowsAffected == 0 {
		log.Println("Exercise not found")
		return &apperrors.AppError{
			Code:    http.StatusNotFound,
			Message: "Exercise not found",
		}
	}

	s.redis.Del(ctx, exerciseCacheKey)
	return nil
}
