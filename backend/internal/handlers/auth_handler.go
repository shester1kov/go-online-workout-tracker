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

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register godoc
// @Summary User registration
// @Description Endpoint for new user registration
// @Tags auth
// @Accept json
// @Produce json
// @Param data body models.UserRegisterRequest true "User data (username, password, email)"
// @Success 201 {object} models.UserResponse "User registered successfully"
// @Failure 400 {object} models.ErrorResponse "Incorrect user data"
// @Failure 400 {object} models.ErrorResponse "Request cancelled"
// @Failure 409 {object} models.ErrorResponse "Email already exists"
// @Failure 409 {object} models.ErrorResponse "Username already exists"
// @Failure 500 {object} models.ErrorResponse "Failed to register user"
// @Failure 504 {object} models.ErrorResponse "Request timeout"
// @Router /register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var req models.UserRegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.JSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.authService.Register(ctx, &req)
	if err != nil {
		log.Println("Registration failed:", err)
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
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// Login godoc
// @Summary User login
// @Description Endpoint for login
// @Tags auth
// @Accept json
// @Produce json
// @Param data body models.UserAuthRequest true "User data (email, password)"
// @Success 200 {object} models.UserResponse "User login successfully"
// @Failure 400 {object} models.ErrorResponse "Incorrect user data"
// @Failure 400 {object} models.ErrorResponse "Invalid password"
// @Failure 400 {object} models.ErrorResponse "Request cancelled"
// @Failure 404 {object} models.ErrorResponse "User not found"
// @Failure 500 {object} models.ErrorResponse "Failed to login user"
// @Failure 504 {object} models.ErrorResponse "Request timeout"
// @Router /login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	log.Println("CORS Headers in main handler:", w.Header())
	ctx := r.Context()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var req models.UserAuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.JSONError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	token, user, err := h.authService.Login(ctx, &req)
	if err != nil {
		log.Println("Login failed:", err)
		var appErr *apperrors.AppError
		if errors.As(err, &appErr) {
			utils.JSONError(w, appErr.Message, appErr.Code)
			return
		}
		utils.JSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    token,
		HttpOnly: true,
		Secure:   false, // true = only HTTPS (possible false on localhost)
		SameSite: http.SameSiteLaxMode,
		Path:     "/api",
		MaxAge:   900,
	})

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

// Logout godoc
// @Summary User logout
// @Description Endpoint for logout
// @Tags auth
// @Success 200
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 403 {object} models.ErrorResponse "Forbidden"
// @Router /logout [post]
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		HttpOnly: true,
		Secure:   false, // true = only HTTPS (possible false on localhost)
		SameSite: http.SameSiteLaxMode,
		Path:     "/api",
		MaxAge:   -1,
	})

	w.WriteHeader(http.StatusOK)
}
