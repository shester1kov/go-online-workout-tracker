package repository

import (
	"backend/internal/models"
	"context"
	"log"

	"github.com/jmoiron/sqlx"
)

type WorkoutExerciseRepository struct {
	db *sqlx.DB
}

func NewWorkoutExerciseRepository(db *sqlx.DB) *WorkoutExerciseRepository {
	return &WorkoutExerciseRepository{db: db}
}

func (r *WorkoutExerciseRepository) AddExerciseToWorkout(ctx context.Context, workoutExercise *models.WorkoutExercise) error {
	query := `INSERT INTO WorkoutExercises (workout_id, exercise_id, sets, reps, weight, notes)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id, created_at`

	err := r.db.QueryRowContext(
		ctx,
		query,
		workoutExercise.WorkoutID,
		workoutExercise.ExerciseID,
		workoutExercise.Sets,
		workoutExercise.Reps,
		workoutExercise.Weight,
		workoutExercise.Notes,
	).Scan(
		&workoutExercise.ID,
		&workoutExercise.CreatedAt,
	)

	if err != nil {
		log.Println("Failed to add exercise to workout:", err)
		return err
	}

	return nil
}

func (r *WorkoutExerciseRepository) GetExercisesByWorkoutID(ctx context.Context, workoutID int) (*[]models.WorkoutExercise, error) {
	query := `SELECT id, workout_id, exercise_id, sets, reps, weight, notes, created_at
	FROM workoutExercises
	WHERE workout_id = $1`

	rows, err := r.db.QueryContext(
		ctx,
		query,
		workoutID,
	)
	if err != nil {
		log.Println("Failed to get exercises:", err)
		return nil, err
	}
	defer rows.Close()

	var exercises []models.WorkoutExercise
	for rows.Next() {
		var workoutExercise models.WorkoutExercise
		err := rows.Scan(
			&workoutExercise.ID,
			&workoutExercise.WorkoutID,
			&workoutExercise.ExerciseID,
			&workoutExercise.Sets,
			&workoutExercise.Reps,
			&workoutExercise.Weight,
			&workoutExercise.Notes,
			&workoutExercise.CreatedAt,
		)
		if err != nil {
			log.Println("Failed to scan workout exercise:", err)
			return nil, err
		}

		exercises = append(exercises, workoutExercise)
	}

	if err := rows.Err(); err != nil {
		log.Println("Rows error:", err)
		return nil, err
	}

	return &exercises, nil
}
