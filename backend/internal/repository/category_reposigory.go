package repository

import (
	"backend/internal/models"
	"context"
	"database/sql"
	"log"
)

type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) CreateCategory(ctx context.Context, category *models.Category) error {
	query := `INSERT INTO Categories (name, slug, description)
	VALUES ($1, $2, $3)
	RETURNING id, created_at, updated_at`

	err := r.db.QueryRowContext(
		ctx,
		query,
		category.Name,
		category.Slug,
		category.Description,
	).Scan(
		&category.ID,
		&category.CreatedAt,
		&category.UpdatedAt,
	)
	if err != nil {
		log.Println("Error creating category:", err)
		return err
	}

	return nil
}

func (r *CategoryRepository) GetCategories(ctx context.Context) (*[]models.Category, error) {
	query := `SELECT id, name, slug, description, created_at, updated_at
	FROM Categories
	WHERE is_active = TRUE`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		log.Println("Error querying categories", err)
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category

	for rows.Next() {
		var category models.Category
		if err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.Slug,
			&category.Description,
			&category.CreatedAt,
			&category.UpdatedAt,
		); err != nil {
			log.Println("Failed to get categories rows:", err)
			return nil, err
		}

		categories = append(categories, category)
	}

	if err = rows.Err(); err != nil {
		log.Println("Error get categories rows", err)
		return nil, err
	}

	return &categories, nil
}

func (r *CategoryRepository) GetCategory(ctx context.Context, id int) (*models.Category, error) {
	query := `SELECT id, name, slug, description, created_at, updated_at
	FROM Categories
	WHERE id = $1
	AND is_active = TRUE`

	var category models.Category

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&category.ID,
		&category.Name,
		&category.Slug,
		&category.Description,
		&category.CreatedAt,
		&category.UpdatedAt,
	)
	if err != nil {
		log.Println("Error get category", err)
		return nil, err
	}

	return &category, nil
}

func (r *CategoryRepository) UpdateCategory(ctx context.Context, category *models.Category) error {
	query := `UPDATE Categories
	SET name = $1, slug = $2, description = $3, updated_at = NOW()
	WHERE id = $4
	AND is_active = TRUE
	RETURNING created_at, updated_at`

	err := r.db.QueryRowContext(
		ctx,
		query,
		category.Name,
		category.Slug,
		category.Description,
		category.ID,
	).Scan(&category.CreatedAt, &category.UpdatedAt)
	if err != nil {
		log.Println("Error update category", err)
		return err
	}

	return nil
}

func (r *CategoryRepository) DeleteCategory(ctx context.Context, id int) (int, error) {
	query := `UPDATE Categories
	SET is_active = FALSE, updated_at = NOW()
	WHERE id = $1
	AND is_active = TRUE`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		log.Println("Error delete category", err)
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("Error result of delete category", err)
		return 0, err
	}

	return int(rowsAffected), nil
}
