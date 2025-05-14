package handlers

import (
	"backend/internal/services"
	"backend/internal/utils"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type NutritionHandler struct {
	service *services.NutritionService
}

func NewNutritionHandler(service *services.NutritionService) *NutritionHandler {
	return &NutritionHandler{service: service}
}

// GetDailyNutrition retrieves daily nutrition data
// @Summary Get daily nutrition entries
// @Description Returns nutrition data for the specified date
// @Tags nutrition
// @Produce json
// @Param date path string trye "Date in YYYY-MM-DD format"
// @Success 200 {array} models.NutritionEntry
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /nutritions{date} [get]
func (h *NutritionHandler) GetDailyNutrition(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	dateStr := chi.URLParam(r, "date")

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		utils.JSONError(w, "Invalid date format", http.StatusBadRequest)
		return
	}

	entries, err := h.service.GetDailyNutrition(ctx, date)
	if err != nil {
		log.Println(err.Error())
		utils.JSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(entries)
}
