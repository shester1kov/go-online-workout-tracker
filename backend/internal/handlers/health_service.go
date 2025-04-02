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
