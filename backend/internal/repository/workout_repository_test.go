package repository

import (
	"backend/internal/models"
	"context"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestCreateWorkout(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewWorkoutRepository(sqlxDB)

	ctx := context.Background()
	date := time.Date(2025, 5, 14, 0, 0, 0, 0, time.UTC)
	w := &models.Workout{UserID: 3, Date: date, Notes: "note"}
	created := time.Now()
	updated := created.Add(10 * time.Minute)

	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO Workouts (user_id, date, notes)
	VALUES ($1, $2, $3)
	RETURNING id, created_at, updated_at, is_active`)).
		WithArgs(w.UserID, w.Date, w.Notes).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "is_active"}).
			AddRow(7, created, updated, true))

	err = repo.CreateWorkout(ctx, w)
	assert.NoError(t, err)
	assert.Equal(t, 7, w.ID)
	assert.WithinDuration(t, created, w.CreatedAt, time.Second)
	assert.WithinDuration(t, updated, w.UpdatedAt, time.Second)
	assert.True(t, w.IsActive)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetWorkoutsByUserID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewWorkoutRepository(sqlxDB)

	ctx := context.Background()
	userID := 4
	t1 := time.Date(2025, 5, 10, 0, 0, 0, 0, time.UTC)
	c1 := time.Now()
	u1 := c1.Add(time.Minute)
	t2 := time.Date(2025, 5, 11, 0, 0, 0, 0, time.UTC)
	c2 := time.Now()
	u2 := c2.Add(2 * time.Minute)

	rows := sqlmock.NewRows([]string{"id", "user_id", "date", "notes", "created_at", "updated_at", "is_active"}).
		AddRow(8, userID, t1, "n1", c1, u1, true).
		AddRow(9, userID, t2, "n2", c2, u2, true)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, user_id, date, notes, created_at, updated_at, is_active
	FROM Workouts
	WHERE is_active = TRUE
	AND user_id = $1`)).
		WithArgs(userID).
		WillReturnRows(rows)

	list, err := repo.GetWorkoutsByUserID(ctx, userID)
	assert.NoError(t, err)
	assert.Len(t, *list, 2)
	assert.Equal(t, 8, (*list)[0].ID)
	assert.Equal(t, t1, (*list)[0].Date)
	assert.Equal(t, "n2", (*list)[1].Notes)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetWorkoutByUserID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewWorkoutRepository(sqlxDB)

	ctx := context.Background()
	userID := 5
	workoutID := 10
	d := time.Date(2025, 5, 12, 0, 0, 0, 0, time.UTC)
	c := time.Now()
	u := c.Add(time.Hour)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, user_id, date, notes, created_at, updated_at, is_active
	FROM Workouts
	WHERE is_active = TRUE
	AND user_id = $1
	AND id = $2`)).
		WithArgs(userID, workoutID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "date", "notes", "created_at", "updated_at", "is_active"}).
			AddRow(workoutID, userID, d, "notes", c, u, true))

	w, err := repo.GetWorkoutByUserID(ctx, userID, workoutID)
	assert.NoError(t, err)
	assert.Equal(t, workoutID, w.ID)
	assert.Equal(t, d, w.Date)
	assert.True(t, w.IsActive)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateWorkoutByUserID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewWorkoutRepository(sqlxDB)

	ctx := context.Background()
	w := &models.Workout{ID: 11, UserID: 6, Date: time.Date(2025, 5, 13, 0, 0, 0, 0, time.UTC), Notes: "updn"}
	created := time.Now()
	updated := created.Add(5 * time.Minute)

	mock.ExpectQuery(regexp.QuoteMeta(`UPDATE Workouts 
	SET date = $1, notes = $2, updated_at = NOW()
	WHERE id = $3
	AND is_active = TRUE
	AND user_id = $4
	RETURNING created_at, updated_at, is_active`)).
		WithArgs(w.Date, w.Notes, w.ID, w.UserID).
		WillReturnRows(sqlmock.NewRows([]string{"created_at", "updated_at", "is_active"}).
			AddRow(created, updated, true))

	err = repo.UpdateWorkoutByUserID(ctx, w)
	assert.NoError(t, err)
	assert.WithinDuration(t, updated, w.UpdatedAt, time.Second)
	assert.True(t, w.IsActive)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteWorkoutByUserID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewWorkoutRepository(sqlxDB)

	ctx := context.Background()
	userID := 7
	workoutID := 12

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE Workouts
	SET is_active = FALSE, updated_at = NOW()
	WHERE id = $1
	AND user_id = $2
	AND is_active = TRUE`)).
		WithArgs(workoutID, userID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	n, err := repo.DeleteWorkoutByUserID(ctx, workoutID, userID)
	assert.NoError(t, err)
	assert.Equal(t, 1, n)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateWorkout_ScanError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewWorkoutRepository(sqlxDB)

	ctx := context.Background()
	w := &models.Workout{UserID: 3, Date: time.Now(), Notes: "note"}

	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO Workouts (user_id, date, notes)
	VALUES ($1, $2, $3)
	RETURNING id, created_at, updated_at, is_active`)).
		WithArgs(w.UserID, w.Date, w.Notes).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "is_active"}).AddRow(nil, nil, nil, nil))

	err := repo.CreateWorkout(ctx, w)
	assert.Error(t, err)
}

func TestGetWorkoutsByUserID_QueryError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewWorkoutRepository(sqlxDB)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, user_id, date, notes, created_at, updated_at, is_active
	FROM Workouts
	WHERE is_active = TRUE
	AND user_id = $1`)).
		WillReturnError(fmt.Errorf("query failed"))

	_, err := repo.GetWorkoutsByUserID(context.Background(), 1)
	assert.Error(t, err)
}

func TestGetWorkoutsByUserID_ScanError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewWorkoutRepository(sqlxDB)

	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "user_id", "date", "notes", "created_at", "updated_at", "is_active"}).
		AddRow(1, 1, nil, "n", now, now, true)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, user_id, date, notes, created_at, updated_at, is_active
	FROM Workouts
	WHERE is_active = TRUE
	AND user_id = $1`)).
		WithArgs(1).
		WillReturnRows(rows)

	_, err := repo.GetWorkoutsByUserID(context.Background(), 1)
	assert.Error(t, err)
}

func TestGetWorkoutsByUserID_RowsErr(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewWorkoutRepository(sqlxDB)

	now := time.Now()

	rows := sqlmock.NewRows([]string{"id", "user_id", "date", "notes", "created_at", "updated_at", "is_active"}).
		AddRow(1, 1, "not-a-date", "note", now, now, true)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, user_id, date, notes, created_at, updated_at, is_active
	FROM Workouts
	WHERE is_active = TRUE
	AND user_id = $1`)).
		WithArgs(1).
		WillReturnRows(rows)

	_, err := repo.GetWorkoutsByUserID(context.Background(), 1)
	assert.Error(t, err)
}

func TestGetWorkoutByUserID_ScanError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewWorkoutRepository(sqlxDB)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, user_id, date, notes, created_at, updated_at, is_active
	FROM Workouts
	WHERE is_active = TRUE
	AND user_id = $1
	AND id = $2`)).
		WithArgs(1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "date", "notes", "created_at", "updated_at", "is_active"}).AddRow(nil, nil, nil, nil, nil, nil, nil))

	_, err := repo.GetWorkoutByUserID(context.Background(), 1, 1)
	assert.Error(t, err)
}

func TestUpdateWorkoutByUserID_ScanError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewWorkoutRepository(sqlxDB)

	w := &models.Workout{ID: 1, UserID: 1, Date: time.Now(), Notes: "n"}
	mock.ExpectQuery(regexp.QuoteMeta(`UPDATE Workouts 
	SET date = $1, notes = $2, updated_at = NOW()
	WHERE id = $3
	AND is_active = TRUE
	AND user_id = $4
	RETURNING created_at, updated_at, is_active`)).
		WithArgs(w.Date, w.Notes, w.ID, w.UserID).
		WillReturnRows(sqlmock.NewRows([]string{"created_at", "updated_at", "is_active"}).AddRow(nil, nil, nil))

	err := repo.UpdateWorkoutByUserID(context.Background(), w)
	assert.Error(t, err)
}

func TestDeleteWorkoutByUserID_ExecError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewWorkoutRepository(sqlxDB)

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE Workouts
	SET is_active = FALSE, updated_at = NOW()
	WHERE id = $1
	AND user_id = $2
	AND is_active = TRUE`)).WillReturnError(fmt.Errorf("exec failed"))

	_, err := repo.DeleteWorkoutByUserID(context.Background(), 1, 1)
	assert.Error(t, err)
}

func TestDeleteWorkoutByUserID_RowsAffectedError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewWorkoutRepository(sqlxDB)

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE Workouts
	SET is_active = FALSE, updated_at = NOW()
	WHERE id = $1
	AND user_id = $2
	AND is_active = TRUE`)).
		WithArgs(1, 1).
		WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("rows affected failed")))

	_, err := repo.DeleteWorkoutByUserID(context.Background(), 1, 1)
	assert.Error(t, err)
}
