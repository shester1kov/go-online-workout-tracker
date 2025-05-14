package repository

import (
	"backend/internal/models"
	"context"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestCreateFood(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewFoodRepository(sqlxDB)

	ctx := context.Background()
	now := time.Now()
	foods := []models.Food{
		{UserID: 1, Date: now, Name: "apple", Quantity: 2, Uint: "pcs", WeightGrams: 150, Calories: 95, Protein: 0.5, Carbs: 25, Fat: 0.3},
		{UserID: 1, Date: now, Name: "banana", Quantity: 1, Uint: "pcs", WeightGrams: 120, Calories: 105, Protein: 1.3, Carbs: 27, Fat: 0.4},
	}

	mock.ExpectBegin()
	mock.ExpectPrepare(regexp.QuoteMeta(`INSERT INTO Foods (user_id, date, name, quantity, unit, weight_grams, calories, protein, carbs, fat)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	RETURNING id`))

	for i := range foods {
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO Foods (user_id, date, name, quantity, unit, weight_grams, calories, protein, carbs, fat)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	RETURNING id`)).
			WithArgs(
				foods[i].UserID,
				foods[i].Date,
				foods[i].Name,
				foods[i].Quantity,
				foods[i].Uint,
				foods[i].WeightGrams,
				foods[i].Calories,
				foods[i].Protein,
				foods[i].Carbs,
				foods[i].Fat,
			).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(i + 1))
	}

	mock.ExpectCommit()

	err = repo.CreateFood(ctx, &foods)
	assert.NoError(t, err)
	assert.Equal(t, 1, foods[0].ID)
	assert.Equal(t, 2, foods[1].ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetFoodByDate(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewFoodRepository(sqlxDB)

	ctx := context.Background()
	userID := 2
	date := time.Date(2025, 5, 15, 0, 0, 0, 0, time.UTC)

	rows := sqlmock.NewRows([]string{"id", "user_id", "date", "name", "quantity", "unit", "weight_grams", "calories", "protein", "carbs", "fat"}).
		AddRow(1, userID, date, "milk", 1, "cup", 244, 150, 8, 12, 8).
		AddRow(2, userID, date, "egg", 2, "pcs", 100, 155, 13, 1.1, 11)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, user_id, date, name, quantity, unit, weight_grams, calories, protein, carbs, fat
	FROM Foods
	WHERE user_id = $1
	AND date::date = $2::date
	AND is_active = TRUE`)).
		WithArgs(userID, date).
		WillReturnRows(rows)

	result, err := repo.GetFoodByDate(ctx, date, userID)
	assert.NoError(t, err)
	assert.Len(t, *result, 2)
	assert.Equal(t, "milk", (*result)[0].Name)
	assert.EqualValues(t, 155, (*result)[1].Calories)
	assert.EqualValues(t, 155, (*result)[1].Calories)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateFood_PrepareError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewFoodRepository(sqlxDB)

	ctx := context.Background()
	foods := []models.Food{{UserID: 1}}

	mock.ExpectBegin()
	mock.ExpectPrepare(regexp.QuoteMeta(`INSERT INTO Foods (user_id, date, name, quantity, unit, weight_grams, calories, protein, carbs, fat)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	RETURNING id`)).
		WillReturnError(fmt.Errorf("prepare failed"))
	mock.ExpectRollback()

	err = repo.CreateFood(ctx, &foods)
	assert.Error(t, err)
	assert.EqualError(t, err, "prepare failed")
}

func TestCreateFood_CommitError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewFoodRepository(sqlxDB)

	ctx := context.Background()
	now := time.Now()
	foods := []models.Food{{UserID: 1, Date: now, Name: "a", Quantity: 1, Uint: "u", WeightGrams: 10, Calories: 10, Protein: 1, Carbs: 1, Fat: 1}}

	mock.ExpectBegin()
	mock.ExpectPrepare(regexp.QuoteMeta(`INSERT INTO Foods (user_id, date, name, quantity, unit, weight_grams, calories, protein, carbs, fat)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	RETURNING id`)).
		WillBeClosed()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO Foods (user_id, date, name, quantity, unit, weight_grams, calories, protein, carbs, fat)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	RETURNING id`)).
		WithArgs(
			foods[0].UserID, foods[0].Date, foods[0].Name, foods[0].Quantity,
			foods[0].Uint, foods[0].WeightGrams, foods[0].Calories,
			foods[0].Protein, foods[0].Carbs, foods[0].Fat,
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit().WillReturnError(fmt.Errorf("commit failed"))

	err = repo.CreateFood(ctx, &foods)
	assert.Error(t, err)
	assert.EqualError(t, err, "commit failed")
}

func TestGetFoodByDate_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewFoodRepository(sqlxDB)

	ctx := context.Background()
	date := time.Now()
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, user_id, date, name, quantity, unit, weight_grams, calories, protein, carbs, fat
	FROM Foods
	WHERE user_id = $1
	AND date::date = $2::date
	AND is_active = TRUE`)).
		WillReturnError(fmt.Errorf("query failed"))

	_, err = repo.GetFoodByDate(ctx, date, 1)
	assert.Error(t, err)
	assert.EqualError(t, err, "query failed")
}

func TestGetFoodByDate_ScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewFoodRepository(sqlxDB)

	ctx := context.Background()
	date := time.Now()
	rows := sqlmock.NewRows([]string{"id", "user_id", "date", "name", "quantity", "unit", "weight_grams", "calories", "protein", "carbs", "fat"}).
		AddRow(1, 1, date, nil, 1, "u", 100, 100, 10, 10, 10)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, user_id, date, name, quantity, unit, weight_grams, calories, protein, carbs, fat
	FROM Foods
	WHERE user_id = $1
	AND date::date = $2::date
	AND is_active = TRUE`)).
		WithArgs(1, date).
		WillReturnRows(rows)

	_, err = repo.GetFoodByDate(ctx, date, 1)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Scan")
}
