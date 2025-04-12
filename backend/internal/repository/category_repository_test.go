package repository

import (
	"backend/internal/models"
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestCreateCategory(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")

	repo := NewCategoryRepository(sqlxDB)

	mock.ExpectQuery(`INSERT INTO Categories`).
		WithArgs("Test Category", "test-category", "This is a test category").
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
			AddRow(1, time.Now(), time.Now()))

	category := &models.Category{
		Name:        "Test Category",
		Slug:        "test-category",
		Description: "This is a test category",
	}

	err = repo.CreateCategory(context.Background(), category)
	assert.NoError(t, err)
	assert.Equal(t, 1, category.ID)
}
