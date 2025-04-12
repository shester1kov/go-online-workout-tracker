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

type WorkoutExerciseHandler struct {
	workoutExerciseService *services.WorkoutExerciseSerivce
}

func NewWorkoutExerciseHandler(workoutExerciseService *services.WorkoutExerciseSerivce) *WorkoutExerciseHandler {
	return &WorkoutExerciseHandler{workoutExerciseService: workoutExerciseService}
}

// AddExerciseToWorkout godoc
// @Summary Add exercise to workout
// @Description Add exercise to workout
// @Tags workouts
// @Accept json
// @Produce json
// @Param id path int true "Workout id"
// @Param workoutExercise body models.WorkoutExerciseRequest true "Exercise data"
// @Success 201 {object} models.WorkoutExerciseResponse "Exercise added to workout"
// @Failure 400 {object} models.ErrorResponse "Invalid request body"
// @Failure 400 {object} models.ErrorResponse "Invalid workout id"
// @Failure 400 {object} models.ErrorResponse "Request cancelled"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 403 {object} models.ErrorResponse "Forbidden"
// @Failure 500 {object} models.ErrorResponse "Failed to add exercise"
// @Failure 504 {object} models.ErrorResponse "Request timeout"
// @Router /workouts/{id}/exercises [post]
func (h *WorkoutExerciseHandler) AddExerciseToWorkout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		log.Println("Incorrect id:", err)
		utils.JSONError(w, "Incorrect id", http.StatusBadRequest)
		return
	}

	var request models.WorkoutExerciseRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Println("Invalid input:", err)
		utils.JSONError(w, "Invalid input", http.StatusBadRequest)
		return
	}

	workoutExercise, err := h.workoutExerciseService.AddExerciseToWorkout(ctx, id, &request)
	if err != nil {
		log.Println("Failed to get workout exercises")
		var appErr *apperrors.AppError
		if errors.As(err, &appErr) {
			utils.JSONError(w, appErr.Message, appErr.Code)
			return
		}
		utils.JSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response := models.WorkoutExerciseResponse{
		ID:         workoutExercise.ID,
		WorkoutID:  workoutExercise.WorkoutID,
		ExerciseID: workoutExercise.ExerciseID,
		Sets:       workoutExercise.Sets,
		Reps:       workoutExercise.Reps,
		Weight:     workoutExercise.Weight,
		Notes:      workoutExercise.Notes,
		CreatedAt:  workoutExercise.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)

}

// GetExercisesByWorkoutID godoc
// @Summary Get exercises by workout id
// @Description Get exercises by workout id
// @Tags workouts
// @Accept json
// @Produce json
// @Param id path int true "Workout id"
// @Success 200 {object} models.WorkoutExerciseResponse "Exercise successfully got"
// @Failure 400 {object} models.ErrorResponse "Invalid workout id"
// @Failure 400 {object} models.ErrorResponse "Request cancelled"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 403 {object} models.ErrorResponse "Forbidden"
// @Failure 500 {object} models.ErrorResponse "Failed to add exercise"
// @Failure 504 {object} models.ErrorResponse "Request timeout"
// @Router /workouts/{id}/exercises [get]
func (h *WorkoutExerciseHandler) GetExercisesByWorkoutID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		log.Println("Incorrect id:", err)
		utils.JSONError(w, "Incorrect id", http.StatusBadRequest)
		return
	}

	workoutExercises, err := h.workoutExerciseService.GetExercisesByWorkoutID(ctx, id)
	if err != nil {
		log.Println("Failed to get exercises by workout id:", err)
		var appErr *apperrors.AppError
		if errors.As(err, &appErr) {
			utils.JSONError(w, appErr.Message, appErr.Code)
			return
		}
		utils.JSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var response []models.WorkoutExerciseResponse
	for _, workoutExercise := range *workoutExercises {
		workoutExerciseResponse := models.WorkoutExerciseResponse{
			ID:         workoutExercise.ID,
			WorkoutID:  workoutExercise.WorkoutID,
			ExerciseID: workoutExercise.ExerciseID,
			Sets:       workoutExercise.Sets,
			Reps:       workoutExercise.Reps,
			Weight:     workoutExercise.Weight,
			Notes:      workoutExercise.Notes,
			CreatedAt:  workoutExercise.CreatedAt,
		}

		response = append(response, workoutExerciseResponse)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
