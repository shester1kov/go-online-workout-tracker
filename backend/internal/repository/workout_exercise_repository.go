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
	query := `SELECT we.id, we.workout_id, we.exercise_id, we.sets, we.reps, we.weight, we.notes, we.created_at,
	e.id, e.name, e.description
	FROM WorkoutExercises we
	INNER JOIN Exercises e ON we.exercise_id = e.id
	WHERE we.workout_id = $1`

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
		var exercise models.WorkoutExerciseItem

		err := rows.Scan(
			&workoutExercise.ID,
			&workoutExercise.WorkoutID,
			&workoutExercise.ExerciseID,
			&workoutExercise.Sets,
			&workoutExercise.Reps,
			&workoutExercise.Weight,
			&workoutExercise.Notes,
			&workoutExercise.CreatedAt,
			&exercise.ID,
			&exercise.Name,
			&exercise.Description,
		)
		if err != nil {
			log.Println("Failed to scan workout exercise:", err)
			return nil, err
		}

		if exercise.ID != 0 {
			workoutExercise.Exercise = &exercise
		}

		exercises = append(exercises, workoutExercise)
	}

	if err := rows.Err(); err != nil {
		log.Println("Rows error:", err)
		return nil, err
	}

	return &exercises, nil
}

func (r *WorkoutExerciseRepository) GetExerciseByWorkoutID(ctx context.Context, workoutID, workoutExerciseID int) (*models.WorkoutExercise, error) {
	query := `SELECT we.id, we.workout_id, we.exercise_id, we.sets, we.reps, we.weight, we.notes, we.created_at,
	e.id, e.name, e.description
	FROM WorkoutExercises we
	INNER JOIN Exercises e ON we.exercise_id = e.id
	WHERE we.workout_id = $1
	AND we.id = $2`

	var workoutExercise models.WorkoutExercise
	var exercise models.WorkoutExerciseItem

	err := r.db.QueryRowContext(
		ctx,
		query,
		workoutID,
		workoutExerciseID,
	).Scan(
		&workoutExercise.ID,
		&workoutExercise.WorkoutID,
		&workoutExercise.ExerciseID,
		&workoutExercise.Sets,
		&workoutExercise.Reps,
		&workoutExercise.Weight,
		&workoutExercise.Notes,
		&workoutExercise.CreatedAt,
		&exercise.ID,
		&exercise.Name,
		&exercise.Description,
	)
	if err != nil {
		log.Println("Failed to scan workout exercise:", err)
		return nil, err
	}

	if exercise.ID != 0 {
		workoutExercise.Exercise = &exercise
	}

	return &workoutExercise, nil
}

func (r *WorkoutExerciseRepository) UpdateExerciseInWorkout(ctx context.Context, workoutExercise *models.WorkoutExercise) error {
	query := `UPDATE WorkoutExercises
	SET exercise_id = $1, sets = $2, reps = $3, weight = $4, notes = $5
	WHERE id = $6
	AND workout_id = $7
	RETURNING created_at`

	err := r.db.QueryRowContext(
		ctx,
		query,
		workoutExercise.ExerciseID,
		workoutExercise.Sets,
		workoutExercise.Reps,
		workoutExercise.Weight,
		workoutExercise.Notes,
		workoutExercise.ID,
		workoutExercise.WorkoutID,
	).Scan(
		&workoutExercise.CreatedAt,
	)

	if err != nil {
		log.Println("Failed to update exercise in workout:", err)
		return err
	}

	return nil
}

func (r *WorkoutExerciseRepository) DeleteExerciseByWorkoutID(ctx context.Context, workoutID, workoutExerciseID int) (int, error) {
	query := `DELETE FROM WorkoutExercises
	WHERE id = $1
	AND workout_id = $2`

	result, err := r.db.ExecContext(
		ctx,
		query,
		workoutExerciseID,
		workoutID,
	)

	if err != nil {
		log.Println("Failed to delete exercise in workout:", err)
		return 0, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("Result error:", err)
		return 0, err
	}

	return int(rowsAffected), nil
}
