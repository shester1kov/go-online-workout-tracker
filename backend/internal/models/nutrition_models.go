package models

import "time"

type NutritionEntry struct {
	ID       int       `json:"id"`
	UserID   int       `json:"user_id"`
	FoodName string    `json:"food_name"`
	Calories float64   `json:"calories"`
	Protein  float64   `json:"protein"`
	Fat      float64   `json:"fat"`
	Carbs    float64   `json:"carbs"`
	Date     time.Time `json:"date"`
}

type NutritionResponse struct {
	FoodName string  `json:"food_name"`
	Calories float64 `json:"calories"`
	Protein  float64 `json:"protein"`
	Fat      float64 `json:"fat"`
	Carbs    float64 `json:"carbs"`
}
