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

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// GetCurrentUser  godoc
// @Summary User profile
// @Description Endpoint for get information about user
// @Tags user
// @Produce json
// @Success 200 {object} models.UserResponse "User profile data"
// @Failure 400 {object} models.ErrorResponse "Request cancelled"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 403 {object} models.ErrorResponse "Forbidden"
// @Failure 404 {object} models.ErrorResponse "User not found"
// @Failure 500 {object} models.ErrorResponse "Failed to login user"
// @Failure 504 {object} models.ErrorResponse "Request timeout"
// @Router /users/me [get]
func (h *UserHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	user, err := h.userService.GetUserByID(ctx)
	if err != nil {
		log.Println("Failed to get user profile:", err)
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

// Add role to user godoc
// @Summary Add role to user
// @Description Endpoint for add role to user
// @Tags user
// @Produce json
// @Param id path int true "User id"
// @Success 200 {object} models.UserResponse "User data"
// @Failure 400 {object} models.ErrorResponse "Request cancelled"
// @Failure 400 {object} models.ErrorResponse "Incorrect id"
// @Failure 400 {object} models.ErrorResponse "Incorrect role id"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 403 {object} models.ErrorResponse "Forbidden"
// @Failure 404 {object} models.ErrorResponse "User not found"
// @Failure 409 {object} models.ErrorResponse "User already has this role"
// @Failure 500 {object} models.ErrorResponse "Failed to login user"
// @Failure 504 {object} models.ErrorResponse "Request timeout"
// @Router /user/{id}/roles [post]
func (h *UserHandler) AddRoleToUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		log.Println("Incorrect id:", err)
		utils.JSONError(w, "Incorrect id", http.StatusBadRequest)
		return
	}

	var req models.AddRoleToUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.JSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	userID := ctx.Value("user_id").(int)
	if userID == id {
		utils.JSONError(w, "You can't change your role", http.StatusForbidden)
		return
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	user, err := h.userService.AddRoleToUser(ctx, id, &req)
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
