package repository

import (
	"backend/internal/models"
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestAddExerciseToWorkout(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewWorkoutExerciseRepository(sqlxDB)

	ctx := context.Background()
	we := &models.WorkoutExercise{WorkoutID: 1, ExerciseID: 2, Sets: 3, Reps: 10, Weight: 50.5, Notes: "note"}
	createdAt := time.Now()

	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO WorkoutExercises (workout_id, exercise_id, sets, reps, weight, notes)
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id, created_at`)).
		WithArgs(we.WorkoutID, we.ExerciseID, we.Sets, we.Reps, we.Weight, we.Notes).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).AddRow(5, createdAt))

	err = repo.AddExerciseToWorkout(ctx, we)
	assert.NoError(t, err)
	assert.Equal(t, 5, we.ID)
	assert.WithinDuration(t, createdAt, we.CreatedAt, time.Second)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetExercisesByWorkoutID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewWorkoutExerciseRepository(sqlxDB)

	ctx := context.Background()
	workoutID := 1
	createdAt := time.Now()

	rows := sqlmock.NewRows([]string{
		"id", "workout_id", "exercise_id", "sets", "reps", "weight", "notes", "created_at",
		"id", "name", "description",
	}).
		AddRow(5, workoutID, 2, 3, 10, 50.5, "note", createdAt, 2, "ex", "desc").
		AddRow(6, workoutID, 3, 4, 8, 40.0, "note2", createdAt, 3, "ex2", "desc2")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT we.id, we.workout_id, we.exercise_id, we.sets, we.reps, we.weight, we.notes, we.created_at,
	e.id, e.name, e.description
	FROM WorkoutExercises we
	INNER JOIN Exercises e ON we.exercise_id = e.id
	WHERE we.workout_id = $1`)).
		WithArgs(workoutID).
		WillReturnRows(rows)

	exs, err := repo.GetExercisesByWorkoutID(ctx, workoutID)
	assert.NoError(t, err)
	assert.Len(t, *exs, 2)
	assert.Equal(t, 5, (*exs)[0].ID)
	assert.Equal(t, "ex", (*exs)[0].Exercise.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetExerciseByWorkoutID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewWorkoutExerciseRepository(sqlxDB)

	ctx := context.Background()
	workoutID := 1
	exID := 5
	createdAt := time.Now()

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT we.id, we.workout_id, we.exercise_id, we.sets, we.reps, we.weight, we.notes, we.created_at,
	e.id, e.name, e.description
	FROM WorkoutExercises we
	INNER JOIN Exercises e ON we.exercise_id = e.id
	WHERE we.workout_id = $1
	AND we.id = $2`)).
		WithArgs(workoutID, exID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "workout_id", "exercise_id", "sets", "reps", "weight", "notes", "created_at", "id", "name", "description"}).
			AddRow(exID, workoutID, 2, 3, 10, 50.5, "note", createdAt, 2, "ex", "desc"))

	we, err := repo.GetExerciseByWorkoutID(ctx, workoutID, exID)
	assert.NoError(t, err)
	assert.Equal(t, exID, we.ID)
	assert.Equal(t, "ex", we.Exercise.Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateExerciseInWorkout(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewWorkoutExerciseRepository(sqlxDB)

	ctx := context.Background()
	we := &models.WorkoutExercise{ID: 5, WorkoutID: 1, ExerciseID: 3, Sets: 4, Reps: 12, Weight: 60.0, Notes: "upd"}
	createdAt := time.Now()

	mock.ExpectQuery(regexp.QuoteMeta(`UPDATE WorkoutExercises
	SET exercise_id = $1, sets = $2, reps = $3, weight = $4, notes = $5
	WHERE id = $6
	AND workout_id = $7
	RETURNING created_at`)).
		WithArgs(we.ExerciseID, we.Sets, we.Reps, we.Weight, we.Notes, we.ID, we.WorkoutID).
		WillReturnRows(sqlmock.NewRows([]string{"created_at"}).AddRow(createdAt))

	err = repo.UpdateExerciseInWorkout(ctx, we)
	assert.NoError(t, err)
	assert.WithinDuration(t, createdAt, we.CreatedAt, time.Second)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteExerciseByWorkoutID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewWorkoutExerciseRepository(sqlxDB)

	ctx := context.Background()
	workoutID := 1
	exID := 5

	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM WorkoutExercises
	WHERE id = $1
	AND workout_id = $2`)).
		WithArgs(exID, workoutID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	n, err := repo.DeleteExerciseByWorkoutID(ctx, workoutID, exID)
	assert.NoError(t, err)
	assert.Equal(t, 1, n)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestWorkoutExerciseRepositoryNegative(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repository := NewWorkoutExerciseRepository(sqlxDB)

	t.Run("AddExerciseToWorkout error", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO WorkoutExercises`)).
			WillReturnError(errors.New("insert error"))

		err := repository.AddExerciseToWorkout(context.Background(), &models.WorkoutExercise{})
		assert.Error(t, err)
	})

	t.Run("GetExercisesByWorkoutID query error", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT we.id, we.workout_id, we.exercise_id`)).
			WillReturnError(errors.New("select error"))

		_, err := repository.GetExercisesByWorkoutID(context.Background(), 1)
		assert.Error(t, err)
	})

	t.Run("GetExercisesByWorkoutID scan error", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{
			"we.id", "we.workout_id", "we.exercise_id", "we.sets", "we.reps", "we.weight", "we.notes", "we.created_at",
			"e.id", "e.name", "e.description",
		}).AddRow("bad", 1, 1, 3, 12, 50.0, "note", time.Now(), 1, "Push-up", "Chest exercise")

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT we.id, we.workout_id, we.exercise_id`)).
			WillReturnRows(rows)

		_, err := repository.GetExercisesByWorkoutID(context.Background(), 1)
		assert.Error(t, err)
	})

	t.Run("GetExerciseByWorkoutID error", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT we.id, we.workout_id, we.exercise_id`)).
			WillReturnError(errors.New("get by id error"))

		_, err := repository.GetExerciseByWorkoutID(context.Background(), 1, 1)
		assert.Error(t, err)
	})

	t.Run("GetExerciseByWorkoutID scan error", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{
			"we.id", "we.workout_id", "we.exercise_id", "we.sets", "we.reps", "we.weight", "we.notes", "we.created_at",
			"e.id", "e.name", "e.description",
		}).AddRow("bad", 1, 1, 3, 12, 50.0, "note", time.Now(), 1, "Push-up", "Chest exercise")

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT we.id, we.workout_id, we.exercise_id`)).
			WithArgs(1, 1).WillReturnRows(rows)

		_, err := repository.GetExerciseByWorkoutID(context.Background(), 1, 1)
		assert.Error(t, err)
	})

	t.Run("UpdateExerciseInWorkout error", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`UPDATE WorkoutExercises`)).
			WillReturnError(errors.New("update error"))

		err := repository.UpdateExerciseInWorkout(context.Background(), &models.WorkoutExercise{})
		assert.Error(t, err)
	})

	t.Run("DeleteExerciseByWorkoutID exec error", func(t *testing.T) {
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM WorkoutExercises`)).
			WillReturnError(errors.New("delete error"))

		_, err := repository.DeleteExerciseByWorkoutID(context.Background(), 1, 1)
		assert.Error(t, err)
	})

	t.Run("DeleteExerciseByWorkoutID rows affected error", func(t *testing.T) {
		result := sqlmock.NewErrorResult(errors.New("rows affected error"))

		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM WorkoutExercises`)).
			WillReturnResult(result)

		_, err := repository.DeleteExerciseByWorkoutID(context.Background(), 1, 1)
		assert.Error(t, err)
	})
}
