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

type CategoryHandler struct {
	categoryService *services.CategoryService
}

func NewCategoryHandler(categoryService *services.CategoryService) *CategoryHandler {
	return &CategoryHandler{categoryService: categoryService}
}

// CreateCategory godoc
// @Summary Create category
// @Description Create new category
// @Tags categories
// @Accept json
// @Produce json
// @Param category body models.CategoryRequest true "Category data"
// @Success 201 {object} models.CategoryResponse "Category created"
// @Failure 400 {object} models.ErrorResponse "Invalid request body"
// @Failure 400 {object} models.ErrorResponse "Request cancelled"
// @Failure 409 {object} models.ErrorResponse "Category already exists"
// @Failure 500 {object} models.ErrorResponse "Failed to save exercise"
// @Failure 504 {object} models.ErrorResponse "Request timeout"
// @Router /categories [post]
func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var req models.CategoryRequest

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

	category, err := h.categoryService.CreateCategory(ctx, &req)
	if err != nil {
		log.Println("Failed to create category:", err)
		var appErr *apperrors.AppError
		if errors.As(err, &appErr) {
			utils.JSONError(w, appErr.Message, appErr.Code)
			return
		}
		utils.JSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response := models.CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Slug:        category.Slug,
		Description: category.Description,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// GetCategories godoc
// @Summary Get all categories
// @Description Get all categories from the database
// @Tags categories
// @Accept json
// @Produce json
// @Success 200 {array} models.CategoryResponse "List of categories"
// @Failure 400 {object} models.ErrorResponse "Request cancelled"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 403 {object} models.ErrorResponse "Forbidden"
// @Failure 404 {object} models.ErrorResponse "Exercises not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Failure 504 {object} models.ErrorResponse "Request timeout"
// @Router /categories [get]
func (h *CategoryHandler) GetCategories(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	categories, err := h.categoryService.GetCategories(ctx)
	if err != nil {
		log.Println("Failed to get categories:", err)
		var appErr *apperrors.AppError
		if errors.As(err, &appErr) {
			utils.JSONError(w, appErr.Message, appErr.Code)
			return
		}
		utils.JSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var response []models.CategoryResponse

	for _, category := range *categories {
		response = append(response, models.CategoryResponse{
			ID:          category.ID,
			Name:        category.Name,
			Slug:        category.Slug,
			Description: category.Description,
			CreatedAt:   category.CreatedAt,
			UpdatedAt:   category.UpdatedAt,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}

// GetExercise godoc
// @Summary Get exercise
// @Description Get exercise by id
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "Category id"
// @Success 200 {object} models.CategoryResponse "Category data"
// @Failure 400 {object} models.ErrorResponse "Invalid id"
// @Failure 400 {object} models.ErrorResponse "Request cancelled"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 403 {object} models.ErrorResponse "Forbidden"
// @Failure 404 {object} models.ErrorResponse "Category not found"
// @Failure 500 {object} models.ErrorResponse "Error data receiving"
// @Failure 504 {object} models.ErrorResponse "Request timeout"
// @Router /categories/{id} [get]
func (h *CategoryHandler) GetCategory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	idStr := chi.URLParam(r, "id")

	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		log.Println(err)
		utils.JSONError(w, "Incorrect id:", http.StatusBadRequest)
		return
	}

	category, err := h.categoryService.GetCategory(ctx, id)
	if err != nil {
		log.Println("Failed to get category:", err)
		var appErr *apperrors.AppError
		if errors.As(err, &appErr) {
			utils.JSONError(w, appErr.Message, appErr.Code)
			return
		}
		utils.JSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response := models.CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Slug:        category.Slug,
		Description: category.Description,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// UpdateExercise godoc
// @Summary Update category
// @Description Update category by id
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "Category id"
// @Param category body models.CategoryRequest true "Category data"
// @Success 200 {object} models.CategoryResponse "Category successfully updated"
// @Failure 400 {object} models.ErrorResponse "Invalid id"
// @Failure 400 {object} models.ErrorResponse "Invalid request body"
// @Failure 400 {object} models.ErrorResponse "Request cancelled"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 403 {object} models.ErrorResponse "Forbidden"
// @Failure 404 {object} models.ErrorResponse "Category not found"
// @Failure 409 {object} models.ErrorResponse "Category already exists"
// @Failure 500 {object} models.ErrorResponse "Failed to update category"
// @Failure 504 {object} models.ErrorResponse "Request timeout"
// @Router /categories/{id} [put]
func (h *CategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
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

	var req models.CategoryRequest

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

	category, err := h.categoryService.UpdateCategory(ctx, id, &req)
	if err != nil {
		log.Println("Failed to update category:", err)
		var appErr *apperrors.AppError
		if errors.As(err, &appErr) {
			utils.JSONError(w, appErr.Message, appErr.Code)
			return
		}
		utils.JSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response := models.CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Slug:        category.Slug,
		Description: category.Description,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// DeleteCategory godoc
// @Summary Delete exercise
// @Description Delete exercise by id
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "Category id"
// @Success 204 "Category successfully deleted"
// @Failure 400 {object} models.ErrorResponse "Invalid id"
// @Failure 400 {object} models.ErrorResponse "Request cancelled"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 403 {object} models.ErrorResponse "Forbidden"
// @Failure 404 {object} models.ErrorResponse "Category not found"
// @Failure 500 {object} models.ErrorResponse "Failed to delete exercise"
// @Failure 504 {object} models.ErrorResponse "Request timeout"
// @Router /categories/{id} [delete]
func (h *CategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
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

	err = h.categoryService.DeleteCategory(ctx, id)
	if err != nil {
		log.Println("Failed to delete category:", err)
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
