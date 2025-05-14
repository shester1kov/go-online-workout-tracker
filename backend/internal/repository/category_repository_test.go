package repository

import (
	"backend/internal/models"
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestGetCategory(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewCategoryRepository(sqlxDB)

	ctx := context.Background()
	now := time.Now()
	id := 3

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, slug, description, created_at, updated_at
	FROM Categories
	WHERE id = $1
	AND is_active = TRUE`)).
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "slug", "description", "created_at", "updated_at"}).
			AddRow(id, "c3", "c-3", "d3", now, now))

	cat, err := repo.GetCategory(ctx, id)
	assert.NoError(t, err)
	assert.Equal(t, id, cat.ID)
	assert.Equal(t, "c3", cat.Name)
	assert.WithinDuration(t, now, cat.UpdatedAt, time.Second)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateCategory(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewCategoryRepository(sqlxDB)

	ctx := context.Background()
	now := time.Now()
	cat := &models.Category{ID: 4, Name: "c4", Slug: "c-4", Description: "d4"}

	mock.ExpectQuery(regexp.QuoteMeta(`UPDATE Categories
	SET name = $1, slug = $2, description = $3, updated_at = NOW()
	WHERE id = $4
	AND is_active = TRUE
	RETURNING created_at, updated_at`)).
		WithArgs(cat.Name, cat.Slug, cat.Description, cat.ID).
		WillReturnRows(sqlmock.NewRows([]string{"created_at", "updated_at"}).
			AddRow(now, now))

	err = repo.UpdateCategory(ctx, cat)
	assert.NoError(t, err)
	assert.WithinDuration(t, now, cat.CreatedAt, time.Second)
	assert.WithinDuration(t, now, cat.UpdatedAt, time.Second)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteCategory(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewCategoryRepository(sqlxDB)

	ctx := context.Background()
	id := 5

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE Categories
	SET is_active = FALSE, updated_at = NOW()
	WHERE id = $1
	AND is_active = TRUE`)).
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(0, 1))

	n, err := repo.DeleteCategory(ctx, id)
	assert.NoError(t, err)
	assert.Equal(t, 1, n)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCategoryRepository_CreateCategory_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewCategoryRepository(sqlxDB)

	ctx := context.Background()
	category := &models.Category{
		Name:        "Test",
		Slug:        "test",
		Description: "desc",
	}

	mock.ExpectQuery(`INSERT INTO Categories`).
		WithArgs(category.Name, category.Slug, category.Description).
		WillReturnError(errors.New("db insert error"))

	err = repo.CreateCategory(ctx, category)
	require.Error(t, err)
	require.Contains(t, err.Error(), "db insert error")

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestCategoryRepository_GetCategories_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewCategoryRepository(sqlxDB)

	ctx := context.Background()

	mock.ExpectQuery(`SELECT id, name, slug, description, created_at, updated_at FROM Categories`).
		WillReturnError(errors.New("db select error"))

	categories, err := repo.GetCategories(ctx)
	require.Error(t, err)
	require.Nil(t, categories)
	require.Contains(t, err.Error(), "db select error")

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestCategoryRepository_GetCategory_ScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewCategoryRepository(sqlxDB)

	ctx := context.Background()
	categoryID := 1

	mock.ExpectQuery(`SELECT id, name, slug, description, created_at, updated_at FROM Categories WHERE id = \$1 AND is_active = TRUE`).
		WithArgs(categoryID).
		WillReturnError(sql.ErrNoRows)

	category, err := repo.GetCategory(ctx, categoryID)
	require.Error(t, err)
	require.Nil(t, category)
	require.Equal(t, sql.ErrNoRows, err)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestCategoryRepository_UpdateCategory_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewCategoryRepository(sqlxDB)

	ctx := context.Background()
	category := &models.Category{
		ID:          1,
		Name:        "Updated Name",
		Slug:        "updated-slug",
		Description: "Updated desc",
	}

	mock.ExpectQuery(`UPDATE Categories`).
		WithArgs(category.Name, category.Slug, category.Description, category.ID).
		WillReturnError(errors.New("db update error"))

	err = repo.UpdateCategory(ctx, category)
	require.Error(t, err)
	require.Contains(t, err.Error(), "db update error")

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestCategoryRepository_DeleteCategory_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewCategoryRepository(sqlxDB)

	ctx := context.Background()
	categoryID := 1

	mock.ExpectExec(`UPDATE Categories`).
		WithArgs(categoryID).
		WillReturnError(errors.New("db exec error"))

	rowsAffected, err := repo.DeleteCategory(ctx, categoryID)
	require.Error(t, err)
	require.Equal(t, 0, rowsAffected)
	require.Contains(t, err.Error(), "db exec error")

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}
