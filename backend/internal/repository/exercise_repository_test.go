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
	"github.com/stretchr/testify/require"
)

func setupMockRepo(t *testing.T) (*ExerciseRepository, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("ошибка при создании sqlmock: %s", err)
	}
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewExerciseRepository(sqlxDB)
	teardown := func() {
		db.Close()
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("ожидания не выполнены: %s", err)
		}
	}
	return repo, mock, teardown
}

func TestCreateExercise(t *testing.T) {
	repo, mock, teardown := setupMockRepo(t)
	defer teardown()

	ctx := context.Background()
	now := time.Now()

	exercise := &models.Exercise{
		Name:        "Приседания",
		Description: "Упражнение для ног",
		CategoryID:  1,
	}

	mock.ExpectQuery(regexp.QuoteMeta(`
		INSERT INTO Exercises (name, description, category_id)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`)).
		WithArgs(exercise.Name, exercise.Description, exercise.CategoryID).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
				AddRow(1, now, now),
		)

	err := repo.CreateExercise(ctx, exercise)
	if err != nil {
		t.Errorf("ошибка при создании упражнения: %s", err)
	}

	if exercise.ID != 1 {
		t.Errorf("ожидался ID 1, получен %d", exercise.ID)
	}
}

func TestGetExercise(t *testing.T) {
	repo, mock, teardown := setupMockRepo(t)
	defer teardown()

	ctx := context.Background()
	now := time.Now()

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT id, name, description, category_id, created_at, updated_at
		FROM Exercises
		WHERE id = $1
		AND is_active = TRUE`,
	)).WithArgs(42).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "category_id", "created_at", "updated_at"}).
			AddRow(42, "жим лёжа", "грудное упражнение", 2, now, now))

	ex, err := repo.GetExercise(ctx, 42)
	if err != nil {
		t.Errorf("ошибка GetExercise: %s", err)
	}

	if ex.ID != 42 {
		t.Errorf("ожидался id 42, получили %d", ex.ID)
	}
	if ex.Name != "жим лёжа" {
		t.Errorf("ожидалось имя 'жим лёжа', получили %s", ex.Name)
	}
}

func TestUpdateExercise(t *testing.T) {
	repo, mock, teardown := setupMockRepo(t)
	defer teardown()

	ctx := context.Background()
	now := time.Now()

	exercise := &models.Exercise{
		ID:          5,
		Name:        "присед",
		Description: "ноги",
		CategoryID:  1,
	}

	mock.ExpectQuery(regexp.QuoteMeta(
		`UPDATE Exercises
		SET name = $1, description = $2, category_id = $3, updated_at = NOW()
		WHERE id = $4
		AND is_active = TRUE
		RETURNING created_at, updated_at`,
	)).WithArgs(exercise.Name, exercise.Description, exercise.CategoryID, exercise.ID).
		WillReturnRows(sqlmock.NewRows([]string{"created_at", "updated_at"}).
			AddRow(now.Add(-time.Hour), now))

	err := repo.UpdateExercise(ctx, exercise)
	if err != nil {
		t.Errorf("ошибка UpdateExercise: %s", err)
	}

	if exercise.UpdatedAt.Before(exercise.CreatedAt) {
		t.Errorf("обновлённое время должно быть позже созданного: created %v, updated %v", exercise.CreatedAt, exercise.UpdatedAt)
	}
}

func TestDeleteExercise(t *testing.T) {
	repo, mock, teardown := setupMockRepo(t)
	defer teardown()

	ctx := context.Background()

	mock.ExpectExec(regexp.QuoteMeta(
		`UPDATE Exercises
		SET is_active = FALSE, updated_at = NOW()
		WHERE id = $1
		AND is_active = TRUE`,
	)).WithArgs(7).
		WillReturnResult(sqlmock.NewResult(0, 1))

	rows, err := repo.DeleteExercise(ctx, 7)
	if err != nil {
		t.Errorf("ошибка DeleteExercise: %s", err)
	}
	if rows != 1 {
		t.Errorf("ожидалось 1 удалённая строка, получили %d", rows)
	}
}

func TestGetExercises(t *testing.T) {
	repo, mock, teardown := setupMockRepo(t)
	defer teardown()

	ctx := context.Background()
	now := time.Now()

	filter := &models.ExerciseFilter{
		CategoryID: func(i int) *int { return &i }(3),
		Search:     func(s string) *string { return &s }("жим"),
		SortBy:     "name",
		SortOrder:  "ASC",
		Limit:      2,
		Offset:     0,
	}

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT COUNT(*) FROM Exercises WHERE is_active = TRUE AND category_id = $1 AND LOWER(name) LIKE $2`,
	)).WithArgs(*filter.CategoryID, "%жим%").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(3))

	mock.ExpectQuery(regexp.QuoteMeta(
		`SELECT id, name, description, category_id, created_at, updated_at FROM Exercises WHERE is_active = TRUE AND category_id = $1 AND LOWER(name) LIKE $2 ORDER BY name ASC LIMIT $3 OFFSET $4`,
	)).WithArgs(*filter.CategoryID, "%жим%", filter.Limit, filter.Offset).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description", "category_id", "created_at", "updated_at"}).
			AddRow(10, "жим лёжа", "грудь", 3, now, now).
			AddRow(11, "жим стоя", "плечи", 3, now, now))

	exs, total, err := repo.GetExercises(ctx, filter)
	if err != nil {
		t.Errorf("ошибка GetExercises: %s", err)
	}

	if total != 3 {
		t.Errorf("ожидалось total=3, получили %d", total)
	}
	if len(*exs) != 2 {
		t.Errorf("ожидалось 2 упражнения, получили %d", len(*exs))
	}
	if (*exs)[0].ID != 10 || (*exs)[1].ID != 11 {
		t.Errorf("ожидались id 10 и 11, получили %d и %d", (*exs)[0].ID, (*exs)[1].ID)
	}
}

func TestExerciseRepository_GetExercises_ErrorCount(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewExerciseRepository(sqlxDB)

	filter := &models.ExerciseFilter{
		Limit:     10,
		Offset:    0,
		SortBy:    "name",
		SortOrder: "ASC",
	}

	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM Exercises WHERE is_active = TRUE`).
		WillReturnError(errors.New("count error"))

	_, _, err = repo.GetExercises(context.Background(), filter)
	require.Error(t, err)
	require.Contains(t, err.Error(), "count error")

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestExerciseRepository_GetExercise_ErrorScan(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewExerciseRepository(sqlxDB)

	mock.ExpectQuery(`SELECT id, name, description, category_id, created_at, updated_at FROM Exercises WHERE id = \$1 AND is_active = TRUE`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "description"}).AddRow(1, "name", "desc")) // не все колонки

	_, err = repo.GetExercise(context.Background(), 1)
	require.Error(t, err)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestExerciseRepository_UpdateExercise_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewExerciseRepository(sqlxDB)

	exercise := &models.Exercise{
		ID:          1,
		Name:        "new name",
		Description: "new desc",
		CategoryID:  2,
	}

	mock.ExpectQuery(`UPDATE Exercises`).
		WithArgs(exercise.Name, exercise.Description, exercise.CategoryID, exercise.ID).
		WillReturnError(errors.New("update error"))

	err = repo.UpdateExercise(context.Background(), exercise)
	require.Error(t, err)
	require.Contains(t, err.Error(), "update error")

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestExerciseRepository_DeleteExercise_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewExerciseRepository(sqlxDB)

	mock.ExpectExec(`UPDATE Exercises`).
		WithArgs(1).
		WillReturnError(errors.New("delete error"))

	rows, err := repo.DeleteExercise(context.Background(), 1)
	require.Error(t, err)
	require.Equal(t, 0, rows)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}
