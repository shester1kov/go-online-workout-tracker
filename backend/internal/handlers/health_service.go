package handlers

import (
	"backend/internal/services"
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type HealthHandler struct {
	healthService *services.HealthService
}

func NewHealthHandler(healthService *services.HealthService) *HealthHandler {
	return &HealthHandler{healthService: healthService}
}

// Check godoc
// @Summary Checking the application's functionality
// @Description Checks the status of the application and its dependencies (database, redis)
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} models.HealthStatus "Application is working correctly"
// @Failure 503 {object} models.HealthStatus "Server or dependency issues"
// @Router /health [get]
func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	status := h.healthService.Check(ctx)

	w.Header().Set("Content-Type", "application/json")
	if !status.IsHealthy() {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
	json.NewEncoder(w).Encode(status)
}
