package models

import "time"

type NutritionixFood struct {
	FoodName           string  `json:"food_name"`
	ServingQty         float64 `json:"serving_qty"`
	ServingUint        string  `json:"serving_unit"`
	ServingWeightGrams float64 `json:"serving_weight_grams"`
	Calories           float64 `json:"nf_calories"`
	Protein            float64 `json:"nf_protein"`
	Carbs              float64 `json:"nf_total_carbohydrate"`
	Fat                float64 `json:"nf_total_fat"`
}

type NutritionixResponse struct {
	Foods []NutritionixFood `json:"foods"`
}

type Food struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	Date        time.Time `json:"date"`
	Name        string    `json:"name"`
	Quantity    float64   `json:"quantity"`
	Uint        string    `json:"unit"`
	WeightGrams float64   `json:"weight_grams"`
	Calories    float64   `json:"calories"`
	Protein     float64   `json:"protein"`
	Carbs       float64   `json:"carbohydrate"`
	Fat         float64   `json:"fat"`
}

type FoodRequestItem struct {
	ProductName string  `json:"product_name"`
	Quantity    float64 `json:"quantity"`
	Unit        string  `json:"unit"`
}

type FoodRequest struct {
	Date  string            `json:"date"`
	Items []FoodRequestItem `json:"items"`
}

type FoodResponseItem struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Quantity    float64 `json:"quantity"`
	Uint        string  `json:"unit"`
	WeightGrams float64 `json:"weight_grams"`
	Calories    float64 `json:"calories"`
	Protein     float64 `json:"protein"`
	Carbs       float64 `json:"carbohydrate"`
	Fat         float64 `json:"fat"`
}

type FoodResponse struct {
	UserID int                `json:"user_id"`
	Date   time.Time          `json:"date"`
	Items  []FoodResponseItem `json:"items"`
}
