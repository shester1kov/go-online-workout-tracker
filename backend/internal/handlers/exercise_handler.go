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
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
)

type ExerciseHandler struct {
	exerciseService *services.ExerciseService
}

func NewExerciseHandler(exerciseService *services.ExerciseService) *ExerciseHandler {
	return &ExerciseHandler{exerciseService: exerciseService}
}

// CreateExercise godoc
// @Summary Create exercise
// @Description Create new exercise
// @Tags exercises
// @Accept json
// @Produce json
// @Param exercise body models.ExerciseRequest true "Exercise data"
// @Success 201 {object} models.ExerciseResponse "Exercise created"
// @Failure 400 {object} models.ErrorResponse "Invalid request body"
// @Failure 400 {object} models.ErrorResponse "Invalid category id"
// @Failure 400 {object} models.ErrorResponse "Request cancelled"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 403 {object} models.ErrorResponse "Forbidden"
// @Failure 409 {object} models.ErrorResponse "Exercise already exists"
// @Failure 500 {object} models.ErrorResponse "Failed to save exercise"
// @Failure 504 {object} models.ErrorResponse "Request timeout"
// @Router /exercises [post]
func (h *ExerciseHandler) CreateExercise(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var req models.ExerciseRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println("Invalid request body:", err)
		utils.JSONError(w, "Invalid request body:", http.StatusBadRequest)
		return
	}

	if len(strings.TrimSpace(req.Name)) < 2 {
		log.Println("Invalid request body:", err)
		utils.JSONError(w, "Invalid request body:", http.StatusBadRequest)
		return
	}

	exercise, err := h.exerciseService.CreateExercise(ctx, &req)
	if err != nil {
		log.Println("Failed to create exercise:", err)
		var appErr *apperrors.AppError
		if errors.As(err, &appErr) {
			utils.JSONError(w, appErr.Message, appErr.Code)
			return
		}
		utils.JSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response := models.ExerciseResponse{
		ID:          exercise.ID,
		Name:        exercise.Name,
		Description: exercise.Description,
		CategoryID:  exercise.CategoryID,
		CreatedAt:   exercise.CreatedAt,
		UpdatedAt:   exercise.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// GetExercises godoc
// @Summary Get all exercises
// @Description Get all exercises from the database
// @Tags exercises
// @Accept json
// @Produce json
// @Success 200 {array} models.ExerciseResponse "List of exercises"
// @Failure 400 {object} models.ErrorResponse "Request cancelled"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 403 {object} models.ErrorResponse "Forbidden"
// @Failure 404 {object} models.ErrorResponse "Exercises not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Failure 504 {object} models.ErrorResponse "Request timeout"
// @Router /exercises [get]
func (h *ExerciseHandler) GetExercises(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	exercises, err := h.exerciseService.GetExercises(ctx)
	if err != nil {
		log.Println("Failed to fecth exercises:", err)
		var appErr *apperrors.AppError
		if errors.As(err, &appErr) {
			utils.JSONError(w, appErr.Message, appErr.Code)
			return
		}
		utils.JSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var response []models.ExerciseResponse

	for _, exercise := range *exercises {
		response = append(response, models.ExerciseResponse{
			ID:          exercise.ID,
			Name:        exercise.Name,
			Description: exercise.Description,
			CategoryID:  exercise.CategoryID,
			CreatedAt:   exercise.CreatedAt,
			UpdatedAt:   exercise.UpdatedAt,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}

// GetExercise godoc
// @Summary Get exercise
// @Description Get exercise by id
// @Tags exercises
// @Accept json
// @Produce json
// @Param id path int true "Exercise id"
// @Success 200 {object} models.ExerciseResponse "Exercise data"
// @Failure 400 {object} models.ErrorResponse "Invalid id"
// @Failure 400 {object} models.ErrorResponse "Request cancelled"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 403 {object} models.ErrorResponse "Forbidden"
// @Failure 404 {object} models.ErrorResponse "Exercise not found"
// @Failure 500 {object} models.ErrorResponse "Error data receiving"
// @Failure 504 {object} models.ErrorResponse "Request timeout"
// @Router /exercises/{id} [get]
func (h *ExerciseHandler) GetExercise(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		log.Println(err)
		utils.JSONError(w, "Incorrect id", http.StatusBadRequest)
		return
	}

	exercise, err := h.exerciseService.GetExercise(ctx, id)
	if err != nil {
		log.Println("Failed to fetch exercise:", err)
		var appErr *apperrors.AppError
		if errors.As(err, &appErr) {
			utils.JSONError(w, appErr.Message, appErr.Code)
			return
		}
		utils.JSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response := models.ExerciseResponse{
		ID:          exercise.ID,
		Name:        exercise.Name,
		Description: exercise.Description,
		CategoryID:  exercise.CategoryID,
		CreatedAt:   exercise.CreatedAt,
		UpdatedAt:   exercise.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// UpdateExercise godoc
// @Summary Update exercise
// @Description Update exercise by id
// @Tags exercises
// @Accept json
// @Produce json
// @Param id path int true "Exercise id"
// @Param exercise body models.ExerciseRequest true "Exercise data"
// @Success 200 {object} models.ExerciseResponse "Exercise successfully updated"
// @Failure 400 {object} models.ErrorResponse "Invalid id"
// @Failure 400 {object} models.ErrorResponse "Invalid request body"
// @Failure 400 {object} models.ErrorResponse "Invalid category id"
// @Failure 400 {object} models.ErrorResponse "Request cancelled"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 403 {object} models.ErrorResponse "Forbidden"
// @Failure 404 {object} models.ErrorResponse "Exercise not found"
// @Failure 409 {object} models.ErrorResponse "Exercise already exists"
// @Failure 500 {object} models.ErrorResponse "Failed to update exercise"
// @Failure 504 {object} models.ErrorResponse "Request timeout"
// @Router /exercises/{id} [put]
func (h *ExerciseHandler) UpdateExercise(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		log.Println(err)
		utils.JSONError(w, "Incorrect id", http.StatusBadRequest)
		return
	}

	var req models.ExerciseRequest

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println("Invalid request body:", err)
		utils.JSONError(w, "Invalid request body:", http.StatusBadRequest)
		return
	}

	if len(strings.TrimSpace(req.Name)) < 2 {
		log.Println("Invalid request body:", err)
		utils.JSONError(w, "Invalid request body:", http.StatusBadRequest)
		return
	}

	exercise, err := h.exerciseService.UpdateExercise(ctx, id, &req)
	if err != nil {
		log.Println("Failed to update exercise:", err)
		var appErr *apperrors.AppError
		if errors.As(err, &appErr) {
			utils.JSONError(w, appErr.Message, appErr.Code)
			return
		}
		utils.JSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response := models.ExerciseResponse{
		ID:          exercise.ID,
		Name:        exercise.Name,
		Description: exercise.Description,
		CategoryID:  exercise.CategoryID,
		CreatedAt:   exercise.CreatedAt,
		UpdatedAt:   exercise.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// DeleteExercise godoc
// @Summary Delete exercise
// @Description Delete exercise by id
// @Tags exercises
// @Accept json
// @Produce json
// @Param id path int true "Exercise id"
// @Success 204 "Exercise successfully deleted"
// @Failure 400 {object} models.ErrorResponse "Invalid id"
// @Failure 400 {object} models.ErrorResponse "Request cancelled"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 403 {object} models.ErrorResponse "Forbidden"
// @Failure 404 {object} models.ErrorResponse "Exercise not found"
// @Failure 500 {object} models.ErrorResponse "Failed to delete exercise"
// @Failure 504 {object} models.ErrorResponse "Request timeout"
// @Router /exercises/{id} [delete]
func (h *ExerciseHandler) DeleteExercise(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		log.Println(err)
		utils.JSONError(w, "Incorrect id", http.StatusBadRequest)
		return
	}

	err = h.exerciseService.DeleteExercise(ctx, id)

	if err != nil {
		log.Println("Failed to delete exercise:", err)
		var appErr *apperrors.AppError
		if errors.As(err, &appErr) {
			utils.JSONError(w, appErr.Message, appErr.Code)
			return
		}
		utils.JSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
