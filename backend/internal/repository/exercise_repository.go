package repository

import (
	"backend/internal/models"
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
)

type ExerciseRepository struct {
	db *sqlx.DB
}

func NewExerciseRepository(db *sqlx.DB) *ExerciseRepository {
	return &ExerciseRepository{db: db}
}

func (r *ExerciseRepository) CreateExercise(ctx context.Context, exercise *models.Exercise) error {
	query := `INSERT INTO Exercises (name, description, category_id)
	VALUES ($1, $2, $3)
	RETURNING id, created_at, updated_at`

	err := r.db.QueryRowContext(
		ctx,
		query,
		exercise.Name,
		exercise.Description,
		exercise.CategoryID,
	).Scan(
		&exercise.ID,
		&exercise.CreatedAt,
		&exercise.UpdatedAt,
	)
	if err != nil {
		log.Println("Failed to create exercise:", err)
		return err
	}

	return nil
}

func (r *ExerciseRepository) GetExercises(ctx context.Context, filter *models.ExerciseFilter) (*[]models.Exercise, int, error) {

	var conditions []string
	var args []interface{}
	paramIndex := 1

	if filter.CategoryID != nil {
		conditions = append(conditions, fmt.Sprintf("category_id = $%d", paramIndex))
		args = append(args, *filter.CategoryID)
		paramIndex++
	}

	if filter.Search != nil {
		conditions = append(conditions, fmt.Sprintf("LOWER(name) LIKE $%d", paramIndex))
		args = append(args, "%"+strings.ToLower(*filter.Search)+"%")
		paramIndex++
	}

	baseQuery := "FROM Exercises WHERE is_active = TRUE"

	if len(conditions) > 0 {
		baseQuery += " AND " + strings.Join(conditions, " AND ")
	}

	countQuery := "SELECT COUNT(*) " + baseQuery
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		log.Println("Failed to get total exercises:", err)
		return nil, 0, err
	}

	query := "SELECT id, name, description, category_id, created_at, updated_at " + baseQuery
	query += fmt.Sprintf(" ORDER BY %s %s", filter.SortBy, filter.SortOrder)
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", paramIndex, paramIndex+1)
	args = append(args, filter.Limit, filter.Offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		log.Println("Failed to get exercises:", err)
		return nil, 0, err
	}
	defer rows.Close()

	var exercises []models.Exercise

	for rows.Next() {
		var exercise models.Exercise

		err := rows.Scan(
			&exercise.ID,
			&exercise.Name,
			&exercise.Description,
			&exercise.CategoryID,
			&exercise.CreatedAt,
			&exercise.UpdatedAt,
		)
		if err != nil {
			log.Println("Failed to get exercises rows:", err)
			return nil, 0, err
		}

		exercises = append(exercises, exercise)
	}

	if err = rows.Err(); err != nil {
		log.Println("Failed to get exercises rows:", err)
		return nil, 0, err
	}

	return &exercises, total, nil
}

func (r *ExerciseRepository) GetExercise(ctx context.Context, id int) (*models.Exercise, error) {
	query := `SELECT id, name, description, category_id, created_at, updated_at
	FROM Exercises
	WHERE id = $1
	AND is_active = TRUE`

	var exercise models.Exercise

	row := r.db.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&exercise.ID,
		&exercise.Name,
		&exercise.Description,
		&exercise.CategoryID,
		&exercise.CreatedAt,
		&exercise.UpdatedAt,
	)
	if err != nil {
		log.Println("Failed to get exercise:", err)
		return nil, err
	}

	return &exercise, nil
}

func (r *ExerciseRepository) UpdateExercise(ctx context.Context, exercise *models.Exercise) error {
	query := `UPDATE Exercises
	SET name = $1, description = $2, category_id = $3, updated_at = NOW()
	WHERE id = $4
	AND is_active = TRUE
	RETURNING created_at, updated_at`

	err := r.db.QueryRowContext(
		ctx,
		query,
		exercise.Name,
		exercise.Description,
		exercise.CategoryID,
		exercise.ID,
	).Scan(&exercise.CreatedAt, &exercise.UpdatedAt)
	if err != nil {
		log.Println("Failed to update exercise:", err)
		return err
	}

	return nil
}

func (r *ExerciseRepository) DeleteExercise(ctx context.Context, id int) (int, error) {

	query := `UPDATE Exercises
	SET is_active = FALSE, updated_at = NOW()
	WHERE id = $1
	AND is_active = TRUE`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		log.Println("Failed to delete exercise", err)
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("Failed to delete exercise result", err)
		return 0, err
	}

	return int(rowsAffected), nil
}
