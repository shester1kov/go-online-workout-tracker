package handlers

import (
	"backend/internal/apperrors"
	"backend/internal/models"
	"backend/internal/services"
	"backend/internal/utils"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

type WorkoutHandler struct {
	workoutSerivce *services.WorkoutSerivce
}

func NewWorkoutHandler(workoutService *services.WorkoutSerivce) *WorkoutHandler {
	return &WorkoutHandler{workoutSerivce: workoutService}
}

// CreateWorkout godoc
// @Summary Create workout
// @Description Create new workout
// @Tags workouts
// @Accept json
// @Produce json
// @Param workout body models.WorkoutRequest true "Workout data"
// @Success 201 {object} models.WorkoutResponse "Workout created"
// @Failure 400 {object} models.ErrorResponse "Invalid request body"
// @Failure 400 {object} models.ErrorResponse "Request cancelled"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 403 {object} models.ErrorResponse "Forbidden"
// @Failure 500 {object} models.ErrorResponse "Failed to create workout"
// @Failure 504 {object} models.ErrorResponse "Request timeout"
// @Router /workouts [post]
func (h *WorkoutHandler) CreateWorkout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var request models.WorkoutRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		utils.JSONError(w, "Invalid input", http.StatusBadRequest)
		return
	}

	workout, err := h.workoutSerivce.CreateWorkout(ctx, &request)
	if err != nil {
		log.Println("Failed to create workout:", err)
		var appErr *apperrors.AppError
		if errors.As(err, &appErr) {
			utils.JSONError(w, appErr.Message, appErr.Code)
			return
		}
		utils.JSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response := models.WorkoutResponse{
		ID:        workout.ID,
		UserID:    workout.UserID,
		Date:      workout.Date,
		Notes:     workout.Notes,
		CreatedAt: workout.CreatedAt,
		UpdatedAt: workout.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)

}

// GetWorkoutsByUserID godoc
// @Summary Get workouts by user id
// @Description Get workouts by user id
// @Tags workouts
// @Accept json
// @Produce json
// @Success 200 {array} models.WorkoutResponse "Workouts got successfully"
// @Failure 400 {object} models.ErrorResponse "Request cancelled"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 403 {object} models.ErrorResponse "Forbidden"
// @Failure 404 {object} models.ErrorResponse "Workouts not found"
// @Failure 500 {object} models.ErrorResponse "Failed to get workouts"
// @Failure 504 {object} models.ErrorResponse "Request timeout"
// @Router /workouts [get]
func (h *WorkoutHandler) GetWorkoutsByUserID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	workouts, err := h.workoutSerivce.GetWorkoutsByUserID(ctx)
	if err != nil {
		log.Println("Failed to get workout")
		var appErr *apperrors.AppError
		if errors.As(err, &appErr) {
			utils.JSONError(w, appErr.Message, appErr.Code)
			return
		}
		utils.JSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var response []models.WorkoutResponse
	for _, workout := range *workouts {
		workoutResponse := models.WorkoutResponse{
			ID:        workout.ID,
			UserID:    workout.UserID,
			Date:      workout.Date,
			Notes:     workout.Notes,
			CreatedAt: workout.CreatedAt,
			UpdatedAt: workout.UpdatedAt,
		}

		response = append(response, workoutResponse)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// GetWorkoutByUserID godoc
// @Summary Get workout by user id
// @Description Get workout by user id
// @Tags workouts
// @Accept json
// @Produce json
// @Param id path int true "Workout id"
// @Success 200 {object} models.WorkoutResponse "Workout got successfully"
// @Failure 400 {object} models.ErrorResponse "Request cancelled"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 403 {object} models.ErrorResponse "Forbidden"
// @Failure 404 {object} models.ErrorResponse "Workout not found"
// @Failure 500 {object} models.ErrorResponse "Failed to get workout"
// @Failure 504 {object} models.ErrorResponse "Request timeout"
// @Router /workouts/{id} [get]
func (h *WorkoutHandler) GetWorkoutByUserID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	workoutIDStr := chi.URLParam(r, "id")

	workoutID, err := strconv.Atoi(workoutIDStr)
	if err != nil || workoutID < 1 {
		log.Println("Incorrect id:", err)
		utils.JSONError(w, "Incorrect id", http.StatusBadRequest)
		return
	}

	workout, err := h.workoutSerivce.GetWorkoutByUserID(ctx, workoutID)
	if err != nil {
		log.Println("Failed to get workout")
		var appErr *apperrors.AppError
		if errors.As(err, &appErr) {
			utils.JSONError(w, appErr.Message, appErr.Code)
			return
		}
		utils.JSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response := models.WorkoutResponse{
		ID:        workout.ID,
		UserID:    workout.UserID,
		Date:      workout.Date,
		Notes:     workout.Notes,
		CreatedAt: workout.CreatedAt,
		UpdatedAt: workout.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
