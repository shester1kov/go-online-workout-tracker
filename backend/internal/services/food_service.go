package services

import (
	"backend/internal/apperrors"
	"backend/internal/clients"
	"backend/internal/models"
	"backend/internal/repository"
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type FoodService struct {
	nutritionixClient *clients.NutritionixClient
	foodRepo          *repository.FoodRepository
}

func NewFoodService(nutritionixClient *clients.NutritionixClient, foodRepo *repository.FoodRepository) *FoodService {
	return &FoodService{
		nutritionixClient: nutritionixClient,
		foodRepo:          foodRepo,
	}
}

func (s *FoodService) AddFood(ctx context.Context, req *models.FoodRequest) (*[]models.Food, error) {
	userID, ok := ctx.Value("user_id").(int)
	if !ok {
		return nil, &apperrors.AppError{
			Code:    http.StatusUnauthorized,
			Message: "Missing userID in context",
		}
	}

	parsedDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, &apperrors.AppError{
			Code:    http.StatusBadRequest,
			Message: "Invalid date format",
		}
	}

	foodData, err := s.nutritionixClient.GetNutritionData(buildNutritionixQuery(req.Items))
	if err != nil {
		log.Println("Nutritionix client error:", err)
		return nil, &apperrors.AppError{
			Code:    http.StatusInternalServerError,
			Message: "Internal server error",
		}
	}

	var foods []models.Food

	for _, f := range foodData.Foods {
		food := models.Food{
			UserID:      userID,
			Date:        parsedDate,
			Name:        f.FoodName,
			Quantity:    f.ServingQty,
			WeightGrams: f.ServingWeightGrams,
			Uint:        f.ServingUint,
			Calories:    f.Calories,
			Protein:     f.Protein,
			Carbs:       f.Carbs,
			Fat:         f.Fat,
		}

		foods = append(foods, food)
	}

	if err := s.foodRepo.CreateFood(ctx, &foods); err != nil {
		log.Println("Food repository create error:", err)
		return nil, &apperrors.AppError{
			Code:    http.StatusInternalServerError,
			Message: "Add food error",
		}
	}

	return &foods, nil
}

func (s *FoodService) GetFoodByDate(ctx context.Context, date string) (*[]models.Food, error) {
	userID, ok := ctx.Value("user_id").(int)
	if !ok {
		return nil, &apperrors.AppError{
			Code:    http.StatusUnauthorized,
			Message: "Missing userID in context",
		}
	}

	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		log.Println("Date parse error:", err)
		return nil, &apperrors.AppError{
			Code:    http.StatusBadRequest,
			Message: "Invalid date format",
		}
	}

	foods, err := s.foodRepo.GetFoodByDate(ctx, parsedDate, userID)
	if err != nil {
		log.Println("Food repository get error:", err)
		return nil, &apperrors.AppError{
			Code:    http.StatusInternalServerError,
			Message: "Get food error",
		}
	}

	return foods, nil
}

func buildNutritionixQuery(items []models.FoodRequestItem) string {
	var parts []string
	for _, item := range items {
		parts = append(parts, fmt.Sprintf("%.0f%s %s", item.Quantity, item.Unit, item.ProductName))
	}

	return strings.Join(parts, " and ")
}
