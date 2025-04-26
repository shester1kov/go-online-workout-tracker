package repository

import (
	"backend/internal/models"
	"context"
	"database/sql"
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

func TestGetCategories(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")

	repo := &CategoryRepository{db: sqlxDB}

	expectedCategories := []models.Category{
		{
			ID:          1,
			Name:        "Category 1",
			Slug:        "category-1",
			Description: "Description 1",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          2,
			Name:        "Category 2",
			Slug:        "category-2",
			Description: "Description 2",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	rows := sqlmock.NewRows([]string{"id", "name", "slug", "description", "created_at", "updated_at"}).
		AddRow(expectedCategories[0].ID, expectedCategories[0].Name, expectedCategories[0].Slug, expectedCategories[0].Description, expectedCategories[0].CreatedAt, expectedCategories[0].UpdatedAt).
		AddRow(expectedCategories[1].ID, expectedCategories[1].Name, expectedCategories[1].Slug, expectedCategories[1].Description, expectedCategories[1].CreatedAt, expectedCategories[1].UpdatedAt)

	mock.ExpectQuery(`SELECT id, name, slug, description, created_at, updated_at FROM Categories WHERE is_active = TRUE`).
		WillReturnRows(rows)

	ctx := context.Background()
	categories, err := repo.GetCategories(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, categories)
	assert.Equal(t, len(expectedCategories), len(*categories))
	assert.Equal(t, expectedCategories[0].Name, (*categories)[0].Name)
	assert.Equal(t, expectedCategories[1].Slug, (*categories)[1].Slug)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestGetCategories_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")

	repo := &CategoryRepository{db: sqlxDB}

	mock.ExpectQuery(`SELECT id, name, slug, description, created_at, updated_at FROM Categories WHERE is_active = TRUE`).
		WillReturnError(sql.ErrConnDone)

	ctx := context.Background()
	categories, err := repo.GetCategories(ctx)

	assert.Error(t, err)
	assert.Nil(t, categories)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}
