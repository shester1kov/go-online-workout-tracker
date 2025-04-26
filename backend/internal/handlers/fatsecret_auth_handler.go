package handlers

import (
	"backend/internal/services"
	"backend/internal/utils"
	"context"
	"log"
	"net/http"
	"time"
)

type FatSecretAuthHandler struct {
	nutritionService *services.NutritionService
}

func NewFatSecretAuthHandler(nutritionService *services.NutritionService) *FatSecretAuthHandler {
	return &FatSecretAuthHandler{nutritionService: nutritionService}
}

// ConnectFatSecret initiates FatSecret OAuth authentication
// @Summary Initiate FatSecret OAuth flow
// @Description Starts the OAuth 1.0 authentication process with FatSecret API
// @Tags fatsecretauthentication
// @Produce json
// @Success 302 {string} string "Redirect to FatSecret authorization page"
// @Failure 500 {object} models.ErrorResponse
// @Router /connect/fatsecret [get]
func (h *FatSecretAuthHandler) ConnectFatSecret(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	authURL, err := h.nutritionService.InitFatSecretAuth(ctx)
	if err != nil {
		log.Println("Error init fat secret auth:", err)
		utils.JSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, authURL, http.StatusFound)
}

// Callback handles FatSecret OAuth callback
// @Summary FatSecret OAuth callback handler
// @Description Handles the callback from FatSecret after user authorization
// @Tags fatsecretauthentication
// @Param oauth_token query string true "OAuth token"
// @Param oauth_verifier query string true "OAuth verifier"
// @Success 302 {string} string "Redirect to profile page with success flag"
// @Failure 500 {object} models.ErrorResponse
// @Router /oauth/fatsecret/callback [get]
func (h *FatSecretAuthHandler) Callback(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("oauth_token")
	verifier := r.URL.Query().Get("oauth_verifier")

	if err := h.nutritionService.CompleteFatSecretAuth(r.Context(), token, verifier); err != nil {
		log.Println("Error:", err)
		utils.JSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "http://localhost:5173/profile", http.StatusFound)
}
