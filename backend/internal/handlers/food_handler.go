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
	"time"

	"github.com/go-chi/chi/v5"
)

type FoodHandler struct {
	foodService *services.FoodService
}

func NewFoodHandler(foodService *services.FoodService) *FoodHandler {
	return &FoodHandler{
		foodService: foodService,
	}
}

// AddFood godoc
// @Summary Add food
// @Description Add user daily food
// @Tags foods
// @Accept json
// @Produce json
// @Param description body models.FoodRequest true "Food description"
// @Success 201 {array} models.FoodResponse "List of foods"
// @Failure 400 {object} models.ErrorResponse "Bad request"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 403 {object} models.ErrorResponse "Forbidden"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Failure 504 {object} models.ErrorResponse "Request timeout"
// @Router /foods [post]
func (h *FoodHandler) AddFood(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var req models.FoodRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println("Incorrect data")
		utils.JSONError(w, "Invalid response data", http.StatusBadRequest)
		return
	}

	if _, err := time.Parse("2006-01-02", req.Date); err != nil {
		log.Println("Invalid date format:", err)
		utils.JSONError(w, "Invalid date format. Use YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	foods, err := h.foodService.AddFood(ctx, &req)
	if err != nil {
		log.Println("Failed to create foods:", err)
		var appErr *apperrors.AppError
		if errors.As(err, &appErr) {
			utils.JSONError(w, appErr.Message, appErr.Code)
			return
		}
		utils.JSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var (
		response          models.FoodResponse
		foodResponseItems []models.FoodResponseItem
	)

	for _, food := range *foods {
		response.UserID = food.UserID
		response.Date = food.Date
		foodResponseItem := models.FoodResponseItem{
			ID:          food.ID,
			Name:        food.Name,
			Quantity:    food.Quantity,
			Uint:        food.Uint,
			WeightGrams: food.WeightGrams,
			Calories:    food.Calories,
			Protein:     food.Protein,
			Carbs:       food.Carbs,
			Fat:         food.Fat,
		}

		foodResponseItems = append(foodResponseItems, foodResponseItem)

	}

	response.Items = foodResponseItems

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)

}

// GetFood godoc
// @Summary Get food
// @Description Get user daily food
// @Tags foods
// @Accept json
// @Produce json
// @Param date path string true "Date"
// @Success 200 {array} models.FoodResponse "List of foods"
// @Failure 400 {object} models.ErrorResponse "Bad request"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 403 {object} models.ErrorResponse "Forbidden"
// @Failure 404 {object} models.ErrorResponse "Not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Failure 504 {object} models.ErrorResponse "Request timeout"
// @Router /foods/{date} [get]
func (h *FoodHandler) GetFood(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	date := chi.URLParam(r, "date")

	if _, err := time.Parse("2006-01-02", date); err != nil {
		log.Println("Invalid date:", err)
		utils.JSONError(w, "Invalid date", http.StatusBadRequest)
		return
	}

	foods, err := h.foodService.GetFoodByDate(ctx, date)
	if err != nil {
		var appErr *apperrors.AppError
		if errors.As(err, &appErr) {
			utils.JSONError(w, appErr.Message, appErr.Code)
			return
		}
		utils.JSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var (
		response          models.FoodResponse
		foodResponseItems []models.FoodResponseItem
	)

	for _, food := range *foods {
		response.UserID = food.UserID
		response.Date = food.Date
		foodResponseItem := models.FoodResponseItem{
			ID:          food.ID,
			Name:        food.Name,
			Quantity:    food.Quantity,
			Uint:        food.Uint,
			WeightGrams: food.WeightGrams,
			Calories:    food.Calories,
			Protein:     food.Protein,
			Carbs:       food.Carbs,
			Fat:         food.Fat,
		}

		foodResponseItems = append(foodResponseItems, foodResponseItem)

	}

	response.Items = foodResponseItems

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}
