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
	"strings"

	"github.com/lib/pq"
)

type WorkoutExerciseSerivce struct {
	workoutRepo         *repository.WorkoutRepository
	workoutExerciseRepo *repository.WorkoutExerciseRepository
	exerciseRepo        *repository.ExerciseRepository
}

func NewWorkoutExerciseService(
	workoutRepo *repository.WorkoutRepository,
	workoutExerciseRepo *repository.WorkoutExerciseRepository,
	exerciseRepo *repository.ExerciseRepository,
) *WorkoutExerciseSerivce {
	return &WorkoutExerciseSerivce{
		workoutRepo:         workoutRepo,
		workoutExerciseRepo: workoutExerciseRepo,
		exerciseRepo:        exerciseRepo,
	}
}

func (s *WorkoutExerciseSerivce) AddExerciseToWorkout(ctx context.Context, workoutID int, request *models.WorkoutExerciseRequest) (*models.WorkoutExercise, error) {
	userID, ok := ctx.Value("user_id").(int)
	if !ok {
		log.Println("Unauthorized")
		return nil, &apperrors.AppError{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
		}
	}

	if workout, err := s.workoutRepo.GetWorkoutByUserID(ctx, userID, workoutID); workout == nil || err != nil {
		log.Println("Workout not found")
		return nil, &apperrors.AppError{
			Code:    http.StatusBadRequest,
			Message: "Workout not found",
		}
	}

	if exercise, err := s.exerciseRepo.GetExercise(ctx, request.ExerciseID); exercise == nil || err != nil {
		log.Println("Incorrect exercise id")
		return nil, &apperrors.AppError{
			Code:    http.StatusBadRequest,
			Message: "Invalid exercise id",
		}
	}

	workoutExercise := models.WorkoutExercise{
		WorkoutID:  workoutID,
		ExerciseID: request.ExerciseID,
		Sets:       request.Sets,
		Reps:       request.Reps,
		Weight:     request.Weight,
		Notes:      request.Notes,
	}

	err := s.workoutExerciseRepo.AddExerciseToWorkout(ctx, &workoutExercise)

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
			if strings.Contains(pgErr.Constraint, "workout_id") {
				return nil, &apperrors.AppError{
					Code:    http.StatusNotFound,
					Message: "Workout not found",
				}
			} else if strings.Contains(pgErr.Constraint, "exercise_id") {
				return nil, &apperrors.AppError{
					Code:    http.StatusBadRequest,
					Message: "Incorrect exercise id",
				}
			}

		default:
			return nil, &apperrors.AppError{
				Code:    http.StatusInternalServerError,
				Message: "Failed to add exercise to workout",
			}
		}
	}

	return &workoutExercise, nil
}

func (s *WorkoutExerciseSerivce) GetExercisesByWorkoutID(ctx context.Context, workoutID int) (*[]models.WorkoutExercise, error) {
	userID, ok := ctx.Value("user_id").(int)
	if !ok {
		log.Println("Unauthorized")
		return nil, &apperrors.AppError{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
		}
	}

	if workout, err := s.workoutRepo.GetWorkoutByUserID(ctx, userID, workoutID); workout == nil || err != nil {
		log.Println("Workout not found")
		return nil, &apperrors.AppError{
			Code:    http.StatusNotFound,
			Message: "Workout not found",
		}
	}

	workoutExercises, err := s.workoutExerciseRepo.GetExercisesByWorkoutID(ctx, workoutID)

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
				Message: "Failed to get workout exercises",
			}
		}
	}

	if workoutExercises == nil || len(*workoutExercises) == 0 {
		log.Println("Exericises not found in workout")
		return nil, &apperrors.AppError{
			Code:    http.StatusNotFound,
			Message: "Exercises not found in workout",
		}
	}

	return workoutExercises, nil
}

func (s *WorkoutExerciseSerivce) GetExerciseByWorkoutID(ctx context.Context, workoutID, workoutExerciseID int) (*models.WorkoutExercise, error) {
	userID, ok := ctx.Value("user_id").(int)
	if !ok {
		log.Println("Unauthorized")
		return nil, &apperrors.AppError{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
		}
	}

	if workout, err := s.workoutRepo.GetWorkoutByUserID(ctx, userID, workoutID); workout == nil || err != nil {
		log.Println("Workout not found")
		return nil, &apperrors.AppError{
			Code:    http.StatusNotFound,
			Message: "Workout not found",
		}
	}

	workoutExercise, err := s.workoutExerciseRepo.GetExerciseByWorkoutID(ctx, workoutID, workoutExerciseID)

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
			log.Println("Exercise not found:", err)
			return nil, &apperrors.AppError{
				Code:    http.StatusNotFound,
				Message: "Exercises not found in workout",
			}

		default:
			log.Println("Unhandled error:", err)
			return nil, &apperrors.AppError{
				Code:    http.StatusInternalServerError,
				Message: "Failed to get workout exercises",
			}
		}
	}

	if workoutExercise == nil {
		log.Println("Exericises not found in workout")
		return nil, &apperrors.AppError{
			Code:    http.StatusNotFound,
			Message: "Exercises not found in workout",
		}
	}

	return workoutExercise, nil
}

func (s *WorkoutExerciseSerivce) UpdateExerciseInWorkout(ctx context.Context, workoutID, workoutExerciseID int, request *models.WorkoutExerciseRequest) (*models.WorkoutExercise, error) {
	userID, ok := ctx.Value("user_id").(int)
	if !ok {
		log.Println("Unauthorized")
		return nil, &apperrors.AppError{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
		}
	}

	if workout, err := s.workoutRepo.GetWorkoutByUserID(ctx, userID, workoutID); workout == nil || err != nil {
		log.Println("Workout not found")
		return nil, &apperrors.AppError{
			Code:    http.StatusBadRequest,
			Message: "Workout not found",
		}
	}

	if exercise, err := s.exerciseRepo.GetExercise(ctx, request.ExerciseID); exercise == nil || err != nil {
		log.Println("Incorrect exercise id")
		return nil, &apperrors.AppError{
			Code:    http.StatusBadRequest,
			Message: "Invalid exercise id",
		}
	}

	workoutExercise := models.WorkoutExercise{
		ID:         workoutExerciseID,
		WorkoutID:  workoutID,
		ExerciseID: request.ExerciseID,
		Sets:       request.Sets,
		Reps:       request.Reps,
		Weight:     request.Weight,
		Notes:      request.Notes,
	}

	err := s.workoutExerciseRepo.UpdateExerciseInWorkout(ctx, &workoutExercise)

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
			log.Println("Exercise not found:", err)
			return nil, &apperrors.AppError{
				Code:    http.StatusNotFound,
				Message: "Exercise not found",
			}

		case errors.As(err, &pgErr) && pgErr.Code == apperrors.PgErrForeignKeyViolation:
			log.Println("Foreign key violation:", pgErr)
			if strings.Contains(pgErr.Constraint, "workout_id") {
				return nil, &apperrors.AppError{
					Code:    http.StatusNotFound,
					Message: "Workout not found",
				}
			} else if strings.Contains(pgErr.Constraint, "exercise_id") {
				return nil, &apperrors.AppError{
					Code:    http.StatusBadRequest,
					Message: "Incorrect exercise id",
				}
			}

		default:
			return nil, &apperrors.AppError{
				Code:    http.StatusInternalServerError,
				Message: "Failed to add exercise to workout",
			}
		}
	}

	return &workoutExercise, nil
}

func (s *WorkoutExerciseSerivce) DeleteExerciseByWorkoutID(ctx context.Context, workoutID, workoutExerciseID int) error {
	userID, ok := ctx.Value("user_id").(int)
	if !ok {
		log.Println("Unauthorized")
		return &apperrors.AppError{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
		}
	}

	if workout, err := s.workoutRepo.GetWorkoutByUserID(ctx, userID, workoutID); workout == nil || err != nil {
		log.Println("Workout not found")
		return &apperrors.AppError{
			Code:    http.StatusNotFound,
			Message: "Workout not found",
		}
	}

	rowsAffected, err := s.workoutExerciseRepo.DeleteExerciseByWorkoutID(ctx, workoutID, workoutExerciseID)

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
				Message: "Failed to get workout exercises",
			}
		}
	}

	if rowsAffected == 0 {
		log.Println("Exericise not found in workout")
		return &apperrors.AppError{
			Code:    http.StatusNotFound,
			Message: "Exercise not found in workout",
		}
	}

	return nil
}
