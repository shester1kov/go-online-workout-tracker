package utils

import (
	"backend/internal/models"
	"encoding/json"
	"net/http"
)

func JSONError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(models.ErrorResponse{
		Message: message,
		Code:    statusCode,
	})
}
