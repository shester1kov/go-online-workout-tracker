package repository

import (
	"backend/internal/models"
	"context"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

type FoodRepository struct {
	db *sqlx.DB
}

func NewFoodRepository(db *sqlx.DB) *FoodRepository {
	return &FoodRepository{
		db: db,
	}
}

func (r *FoodRepository) CreateFood(ctx context.Context, foods *[]models.Food) error {
	tx, err := r.db.Begin()
	if err != nil {
		log.Println("Transaction begin error:", err)
		return err
	}

	query := `INSERT INTO Foods (user_id, date, name, quantity, unit, weight_grams, calories, protein, carbs, fat)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	RETURNING id`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		tx.Rollback()
		log.Println("Estimate transaction error:", err)
		return err
	}
	defer stmt.Close()

	for i, food := range *foods {
		err = stmt.QueryRowContext(
			ctx,
			food.UserID,
			food.Date,
			food.Name,
			food.Quantity,
			food.Uint,
			food.WeightGrams,
			food.Calories,
			food.Protein,
			food.Carbs,
			food.Fat,
		).Scan(&(*foods)[i].ID)
		if err != nil {
			tx.Rollback()
			log.Println("Execution error:", err)
			return err
		}

	}

	if err := tx.Commit(); err != nil {
		log.Println("Commit error:", err)
		return err
	}

	return nil
}

func (r *FoodRepository) GetFoodByDate(ctx context.Context, date time.Time, userID int) (*[]models.Food, error) {
	query := `SELECT id, user_id, date, name, quantity, unit, weight_grams, calories, protein, carbs, fat
	FROM Foods
	WHERE user_id = $1
	AND date::date = $2::date
	AND is_active = TRUE`

	rows, err := r.db.QueryContext(
		ctx,
		query,
		userID,
		date,
	)

	if err != nil {
		log.Println("Query error:", err)
		return nil, err
	}

	var foods []models.Food

	for rows.Next() {
		var food models.Food

		if err := rows.Scan(
			&food.ID,
			&food.UserID,
			&food.Date,
			&food.Name,
			&food.Quantity,
			&food.Uint,
			&food.WeightGrams,
			&food.Calories,
			&food.Protein,
			&food.Carbs,
			&food.Fat,
		); err != nil {
			log.Println("Error scan rows:", err)
			return nil, err
		}

		foods = append(foods, food)
	}

	if err := rows.Err(); err != nil {
		log.Println("Rows err:", err)
		return nil, err
	}

	return &foods, nil
}
