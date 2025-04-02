package services

import (
	"backend/internal/apperrors"
	"backend/internal/models"
	"backend/internal/repository"
	"backend/internal/utils"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

const categoryCacheKey = "categories"

type CategoryService struct {
	categoryRepo *repository.CategoryRepository
	redis        *redis.Client
}

func NewCategoryService(categoryRepo *repository.CategoryRepository, redis *redis.Client) *CategoryService {
	return &CategoryService{
		categoryRepo: categoryRepo,
		redis:        redis,
	}
}

func (s *CategoryService) CreateCategory(ctx context.Context, req *models.CategoryRequest) (*models.Category, error) {
	category := &models.Category{
		Name:        req.Name,
		Slug:        utils.GenerateSlug(req.Name),
		Description: req.Description,
	}

	err := s.categoryRepo.CreateCategory(ctx, category)
	if err != nil {
		var pgErr *pq.Error
		switch {
		case errors.Is(err, context.Canceled):
			log.Println("Request cancelled:", err)
			return nil, &apperrors.AppError{
				Code:    http.StatusBadRequest,
				Message: "Request cancelled",
			}
		case errors.Is(err, context.DeadlineExceeded):
			log.Println("Deadline exceeded:", err)
			return nil, &apperrors.AppError{
				Code:    http.StatusGatewayTimeout,
				Message: "Request timeout",
			}
		case errors.As(err, &pgErr) && pgErr.Code == apperrors.PgErrUniqueViolation:
			log.Println("Unique violation:", pgErr)
			return nil, &apperrors.AppError{
				Code:    http.StatusConflict,
				Message: "Category already exists",
			}
		default:
			log.Println("Unhandled error:", err)
			return nil, &apperrors.AppError{
				Code:    http.StatusInternalServerError,
				Message: "Failed to create category",
			}
		}
	}
	s.redis.Del(ctx, categoryCacheKey)
	return category, nil
}

func (s *CategoryService) GetCategories(ctx context.Context) (*[]models.Category, error) {
	val, err := s.redis.Get(ctx, categoryCacheKey).Result()
	if err == redis.Nil {
		log.Println("cache not found, getting data from DB")
	} else if err != nil {
		log.Println("error getting data from Redis:", err)
	} else {
		var categories []models.Category

		if err := json.Unmarshal([]byte(val), &categories); err == nil {
			log.Println("data retrieved from cache")
			return &categories, nil
		}

		log.Println("error deserializing data from cache:", err)
	}

	categories, err := s.categoryRepo.GetCategories(ctx)
	if err != nil {
		switch {
		case errors.Is(err, context.Canceled):
			log.Println("Request cancelled:", err)
			return nil, &apperrors.AppError{
				Code:    http.StatusBadRequest,
				Message: "Request cancelled",
			}
		case errors.Is(err, context.DeadlineExceeded):
			log.Println("Deadline exceeded:", err)
			return nil, &apperrors.AppError{
				Code:    http.StatusGatewayTimeout,
				Message: "Request timeout",
			}
		default:
			log.Println("Unhandled error:", err)
			return nil, &apperrors.AppError{
				Code:    http.StatusInternalServerError,
				Message: "Failed to get categories",
			}
		}
	}

	if categories == nil || len(*categories) == 0 {
		log.Println("Categories not found:", err)
		return nil, &apperrors.AppError{
			Code:    http.StatusNotFound,
			Message: "Categories not found",
		}
	}

	data, err := json.Marshal(categories)
	if err != nil {
		log.Println("cache serialization error", err)
	} else {
		s.redis.Set(ctx, categoryCacheKey, data, 10*time.Minute)
		log.Println("data written to cache")
	}

	return categories, nil
}

func (s *CategoryService) GetCategory(ctx context.Context, id int) (*models.Category, error) {

	category, err := s.categoryRepo.GetCategory(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, context.Canceled):
			log.Println("Request cancelled:", err)
			return nil, &apperrors.AppError{
				Code:    http.StatusBadRequest,
				Message: "Request cancelled",
			}
		case errors.Is(err, context.DeadlineExceeded):
			log.Println("Deadline exceeded:", err)
			return nil, &apperrors.AppError{
				Code:    http.StatusGatewayTimeout,
				Message: "Request timeout",
			}
		case errors.Is(err, sql.ErrNoRows):
			log.Println("Category not found:", err)
			return nil, &apperrors.AppError{
				Code:    http.StatusNotFound,
				Message: "Category not found",
			}
		default:
			log.Println("Unhandled error:", err)
			return nil, &apperrors.AppError{
				Code:    http.StatusInternalServerError,
				Message: "Failed to get category",
			}
		}
	}

	return category, nil
}

func (s *CategoryService) UpdateCategory(ctx context.Context, id int, req *models.CategoryRequest) (*models.Category, error) {
	category := &models.Category{
		ID:          id,
		Name:        req.Name,
		Slug:        utils.GenerateSlug(req.Name),
		Description: req.Description,
	}

	err := s.categoryRepo.UpdateCategory(ctx, category)
	if err != nil {
		var pgErr *pq.Error
		switch {
		case errors.Is(err, context.Canceled):
			log.Println("Request cancelled:", err)
			return nil, &apperrors.AppError{
				Code:    http.StatusBadRequest,
				Message: "Request cancelled",
			}

		case errors.Is(err, context.DeadlineExceeded):
			log.Println("Deadline exceeded:", err)
			return nil, &apperrors.AppError{
				Code:    http.StatusGatewayTimeout,
				Message: "Request timeout",
			}

		case errors.Is(err, sql.ErrNoRows):
			log.Println("Category not found:", err)
			return nil, &apperrors.AppError{
				Code:    http.StatusNotFound,
				Message: "Category not found",
			}

		case errors.As(err, &pgErr) && pgErr.Code == apperrors.PgErrUniqueViolation:
			log.Println("Unique violation:", pgErr)
			return nil, &apperrors.AppError{
				Code:    http.StatusConflict,
				Message: "Category already exists",
			}

		default:
			log.Println("Unhandled error:", err)
			return nil, &apperrors.AppError{
				Code:    http.StatusInternalServerError,
				Message: "Failed to update category",
			}
		}
	}

	s.redis.Del(ctx, categoryCacheKey)
	return category, nil
}

func (s *CategoryService) DeleteCategory(ctx context.Context, id int) error {
	rowsAffected, err := s.categoryRepo.DeleteCategory(ctx, id)

	if err != nil {
		switch {
		case errors.Is(err, context.Canceled):
			log.Println("Request cancelled:", err)
			return &apperrors.AppError{
				Code:    http.StatusBadRequest,
				Message: "Request cancelled",
			}
		case errors.Is(err, context.DeadlineExceeded):
			log.Println("Deadline exceeded:", err)
			return &apperrors.AppError{
				Code:    http.StatusGatewayTimeout,
				Message: "Request timeout",
			}
		default:
			log.Println("Unhandled error:", err)
			return &apperrors.AppError{
				Code:    http.StatusInternalServerError,
				Message: "Failed to delete category",
			}
		}
	}

	if rowsAffected == 0 {
		return &apperrors.AppError{
			Code:    http.StatusNotFound,
			Message: "Category not found",
		}
	}

	s.redis.Del(ctx, categoryCacheKey)
	return nil
}
