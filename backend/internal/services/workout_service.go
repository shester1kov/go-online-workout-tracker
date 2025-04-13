package services

import (
	"backend/internal/apperrors"
	"backend/internal/models"
	"backend/internal/repository"
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/lib/pq"
)

type WorkoutSerivce struct {
	workoutRepo *repository.WorkoutRepository
}

func NewWorkoutService(workoutRepo *repository.WorkoutRepository) *WorkoutSerivce {
	return &WorkoutSerivce{
		workoutRepo: workoutRepo,
	}
}

func (s *WorkoutSerivce) CreateWorkout(ctx context.Context, req *models.WorkoutRequest) (*models.Workout, error) {
	userID, ok := ctx.Value("user_id").(int)
	if !ok {
		return nil, &apperrors.AppError{
			Code:    http.StatusUnauthorized,
			Message: "Missing userID in context",
		}
	}

	workout := &models.Workout{
		UserID: userID,
		Date:   req.Date,
		Notes:  req.Notes,
	}

	err := s.workoutRepo.CreateWorkout(ctx, workout)
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

		case errors.As(err, &pgErr) && pgErr.Code == apperrors.PgErrForeignKeyViolation:
			log.Println("Foreign key violation:", pgErr)
			return nil, &apperrors.AppError{
				Code:    http.StatusBadRequest,
				Message: "Incorrect user id",
			}

		default:
			log.Println("Unhandled error:", err)
			return nil, &apperrors.AppError{
				Code:    http.StatusInternalServerError,
				Message: "Failed to create workout",
			}
		}
	}

	return workout, nil
}

func (s *WorkoutSerivce) GetWorkoutsByUserID(ctx context.Context) (*[]models.Workout, error) {
	userID, ok := ctx.Value("user_id").(int)
	if !ok {
		log.Println("Unauthorized")
		return nil, &apperrors.AppError{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
		}
	}

	workouts, err := s.workoutRepo.GetWorkoutsByUserID(ctx, userID)
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
				Message: "Failed to get workouts",
			}
		}
	}

	if workouts == nil || len(*workouts) == 0 {
		log.Println("Workouts not found")
		return nil, &apperrors.AppError{
			Code:    http.StatusNotFound,
			Message: "Workouts not found",
		}
	}

	return workouts, nil
}

func (s *WorkoutSerivce) GetWorkoutByUserID(ctx context.Context, workoutID int) (*models.Workout, error) {
	userID, ok := ctx.Value("user_id").(int)
	if !ok {
		log.Println("Unauthorized")
		return nil, &apperrors.AppError{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
		}
	}

	workout, err := s.workoutRepo.GetWorkoutByUserID(ctx, userID, workoutID)
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

		case errors.Is(err, sql.ErrNoRows):
			log.Println("Workout not found")
			return nil, &apperrors.AppError{
				Code:    http.StatusNotFound,
				Message: "Workout not found",
			}

		default:
			log.Println("Unhandled error:", err)
			return nil, &apperrors.AppError{
				Code:    http.StatusInternalServerError,
				Message: "Failed to get workout",
			}
		}
	}

	if workout == nil {
		log.Println("Workout not found")
		return nil, &apperrors.AppError{
			Code:    http.StatusNotFound,
			Message: "Workout not found",
		}
	}

	return workout, nil
}

func (s *WorkoutSerivce) UpdateWorkoutByUserID(ctx context.Context, workoutID int, req *models.WorkoutRequest) (*models.Workout, error) {
	userID, ok := ctx.Value("user_id").(int)
	if !ok {
		return nil, &apperrors.AppError{
			Code:    http.StatusUnauthorized,
			Message: "Missing userID in context",
		}
	}

	workout := &models.Workout{
		ID:     workoutID,
		UserID: userID,
		Date:   req.Date,
		Notes:  req.Notes,
	}

	err := s.workoutRepo.UpdateWorkoutByUserID(ctx, workout)
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
			log.Println("Workout not found:", err)
			return nil, &apperrors.AppError{
				Code:    http.StatusNotFound,
				Message: "Workout not found",
			}

		case errors.As(err, &pgErr) && pgErr.Code == apperrors.PgErrForeignKeyViolation:
			log.Println("Foreign key violation:", pgErr)
			return nil, &apperrors.AppError{
				Code:    http.StatusBadRequest,
				Message: "Incorrect user id",
			}

		default:
			log.Println("Unhandled error:", err)
			return nil, &apperrors.AppError{
				Code:    http.StatusInternalServerError,
				Message: "Failed to update workout",
			}
		}
	}

	return workout, nil
}

func (s *WorkoutSerivce) DeleteWorkoutByUserID(ctx context.Context, id int) error {
	userID, ok := ctx.Value("user_id").(int)
	if !ok {
		return &apperrors.AppError{
			Code:    http.StatusUnauthorized,
			Message: "Missing userID in context",
		}
	}

	rowsAffected, err := s.workoutRepo.DeleteWorkoutByUserID(ctx, id, userID)

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
				Message: "Failed to delete exercise",
			}
		}
	}

	if rowsAffected == 0 {
		log.Println("Exercise not found")
		return &apperrors.AppError{
			Code:    http.StatusNotFound,
			Message: "Workout not found",
		}
	}

	return nil
}
