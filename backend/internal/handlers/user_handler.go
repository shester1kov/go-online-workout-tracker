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
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// Profile godoc
// @Summary User profile
// @Description Endpoint for get user profile
// @Tags user
// @Produce json
// @Success 200 {object} models.UserResponse "User profile data"
// @Failure 400 {object} models.ErrorResponse "Request cancelled"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 403 {object} models.ErrorResponse "Forbidden"
// @Failure 404 {object} models.ErrorResponse "User not found"
// @Failure 500 {object} models.ErrorResponse "Failed to login user"
// @Failure 504 {object} models.ErrorResponse "Request timeout"
// @Router /profile [get]
func (h *UserHandler) Profile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := ctx.Value("user_id").(int)
	if !ok {
		utils.JSONError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := h.userService.GetUserByID(ctx, userID)
	if err != nil {
		utils.JSONError(w, "User not found", http.StatusNotFound)
		return
	}

	response := models.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *UserHandler) AddRoleToUser(w http.ResponseWriter, r *http.Request) {
	var req models.AddRoleToUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.JSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	claims := ctx.Value("user_claims").(*models.Claims)
	if claims.UserID == req.UserID {
		utils.JSONError(w, "You can't change your role", http.StatusForbidden)
		return
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	user, err := h.userService.AddRoleToUser(ctx, &req)
	if err != nil {
		log.Println("Failed to add role:", err)
		var appErr *apperrors.AppError
		if errors.As(err, &appErr) {
			utils.JSONError(w, appErr.Message, appErr.Code)
			return
		}
		utils.JSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response := models.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
