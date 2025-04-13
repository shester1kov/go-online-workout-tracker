package repository

import (
	"backend/internal/models"
	"context"
	"log"

	"github.com/jmoiron/sqlx"
)

type WorkoutRepository struct {
	db *sqlx.DB
}

func NewWorkoutRepository(db *sqlx.DB) *WorkoutRepository {
	return &WorkoutRepository{db: db}
}

func (r *WorkoutRepository) CreateWorkout(ctx context.Context, workout *models.Workout) error {
	query := `INSERT INTO Workouts (user_id, date, notes)
	VALUES ($1, $2, $3)
	RETURNING id, created_at, updated_at, is_active`

	err := r.db.QueryRowContext(
		ctx,
		query,
		workout.UserID,
		workout.Date,
		workout.Notes,
	).Scan(
		&workout.ID,
		&workout.CreatedAt,
		&workout.UpdatedAt,
		&workout.IsActive,
	)

	if err != nil {
		log.Println("Failed to create workout:", err)
		return err
	}

	return nil
}

func (r *WorkoutRepository) GetWorkoutsByUserID(ctx context.Context, userID int) (*[]models.Workout, error) {
	query := `SELECT id, user_id, date, notes, created_at, updated_at, is_active
	FROM Workouts
	WHERE is_active = TRUE
	AND user_id = $1`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		log.Println("Failed to get workouts:", err)
		return nil, err
	}
	defer rows.Close()

	var workouts []models.Workout
	for rows.Next() {
		var workout models.Workout
		err := rows.Scan(
			&workout.ID,
			&workout.UserID,
			&workout.Date,
			&workout.Notes,
			&workout.CreatedAt,
			&workout.UpdatedAt,
			&workout.IsActive,
		)
		if err != nil {
			log.Println("Failed to scan workout:", err)
			return nil, err
		}

		workouts = append(workouts, workout)
	}

	if err := rows.Err(); err != nil {
		log.Println("Rows error:", err)
		return nil, err
	}

	return &workouts, nil
}

func (r *WorkoutRepository) GetWorkoutByUserID(ctx context.Context, userID int, workoutID int) (*models.Workout, error) {
	query := `SELECT id, user_id, date, notes, created_at, updated_at, is_active
	FROM Workouts
	WHERE is_active = TRUE
	AND user_id = $1
	AND id = $2`

	var workout models.Workout

	err := r.db.QueryRowContext(
		ctx,
		query,
		userID,
		workoutID,
	).Scan(
		&workout.ID,
		&workout.UserID,
		&workout.Date,
		&workout.Notes,
		&workout.CreatedAt,
		&workout.UpdatedAt,
		&workout.IsActive,
	)

	if err != nil {
		log.Println("Failed to get workout")
		return nil, err
	}

	return &workout, nil
}

func (r *WorkoutRepository) UpdateWorkoutByUserID(ctx context.Context, workout *models.Workout) error {
	query := `UPDATE Workouts 
	SET date = $1, notes = $2, updated_at = NOW()
	WHERE id = $3
	AND is_active = TRUE
	AND user_id = $4
	RETURNING created_at, updated_at, is_active`

	err := r.db.QueryRowContext(
		ctx,
		query,
		workout.Date,
		workout.Notes,
		workout.ID,
		workout.UserID,
	).Scan(
		&workout.CreatedAt,
		&workout.UpdatedAt,
		&workout.IsActive,
	)

	if err != nil {
		log.Println("Failed to create workout:", err)
		return err
	}

	return nil
}

func (r *WorkoutRepository) DeleteWorkoutByUserID(ctx context.Context, id, userID int) (int, error) {

	query := `UPDATE Workouts
	SET is_active = FALSE, updated_at = NOW()
	WHERE id = $1
	AND user_id = $2
	AND is_active = TRUE`

	result, err := r.db.ExecContext(ctx, query, id, userID)
	if err != nil {
		log.Println("Failed to workout exercise", err)
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("Failed to delete workout result", err)
		return 0, err
	}

	return int(rowsAffected), nil
}
