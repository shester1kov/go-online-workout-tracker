package repository

import (
	"backend/internal/models"
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSaveTempAuth(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewFatSecretAuthRepository(sqlxDB)

	ctx := context.Background()
	auth := &models.TempAuth{UserID: 1, RequestToken: "rt", RequestSecret: "rs"}

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO TempFatsecretAuth (user_id, request_token, request_secret)
         VALUES ($1, $2, $3)
         ON CONFLICT (user_id) DO UPDATE
         SET request_token = $2, request_secret = $3, created_at = NOW()`)).
		WithArgs(auth.UserID, auth.RequestToken, auth.RequestSecret).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.SaveTempAuth(ctx, auth)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetTempAuth(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewFatSecretAuthRepository(sqlxDB)

	ctx := context.Background()
	token := "rt"
	expected := &models.TempAuth{UserID: 2, RequestToken: token, RequestSecret: "rs2"}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT user_id, request_token, request_secret
	FROM TempFatsecretAuth
    WHERE request_token = $1`)).
		WithArgs(token).
		WillReturnRows(sqlmock.NewRows([]string{"user_id", "request_token", "request_secret"}).
			AddRow(expected.UserID, expected.RequestToken, expected.RequestSecret))

	auth, err := repo.GetTempAuth(ctx, token)
	assert.NoError(t, err)
	assert.Equal(t, expected, auth)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSaveAuth(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewFatSecretAuthRepository(sqlxDB)

	ctx := context.Background()
	auth := &models.FatSecretAuth{UserID: 3, AccessToken: "at", AccessSecret: "as"}

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO FatsecretAuth (user_id, access_token, access_secret)
    VALUES ($1, $2, $3)
    ON CONFLICT (user_id) DO UPDATE
    SET access_token = $2, access_secret = $3, updated_at = NOW()`)).
		WithArgs(auth.UserID, auth.AccessToken, auth.AccessSecret).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.SaveAuth(ctx, auth)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAuth(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewFatSecretAuthRepository(sqlxDB)

	ctx := context.Background()
	userID := 4
	expected := &models.FatSecretAuth{UserID: userID, AccessToken: "at4", AccessSecret: "as4"}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT user_id, access_token, access_secret
	FROM FatsecretAuth
	WHERE user_id = $1`)).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{"user_id", "access_token", "access_secret"}).
			AddRow(expected.UserID, expected.AccessToken, expected.AccessSecret))

	auth, err := repo.GetAuth(ctx, userID)
	assert.NoError(t, err)
	assert.Equal(t, expected, auth)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFatSecretAuthRepository_SaveTempAuth_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewFatSecretAuthRepository(sqlxDB)

	ctx := context.Background()
	auth := &models.TempAuth{
		UserID:        1,
		RequestToken:  "token",
		RequestSecret: "secret",
	}

	mock.ExpectExec(`INSERT INTO TempFatsecretAuth`).
		WithArgs(auth.UserID, auth.RequestToken, auth.RequestSecret).
		WillReturnError(errors.New("db exec error"))

	err = repo.SaveTempAuth(ctx, auth)
	require.Error(t, err)
	require.Contains(t, err.Error(), "db exec error")

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestFatSecretAuthRepository_GetTempAuth_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewFatSecretAuthRepository(sqlxDB)

	ctx := context.Background()
	token := "token"

	mock.ExpectQuery(`SELECT user_id, request_token, request_secret FROM TempFatsecretAuth WHERE request_token = \$1`).
		WithArgs(token).
		WillReturnError(sql.ErrNoRows)

	auth, err := repo.GetTempAuth(ctx, token)
	require.Error(t, err)
	require.Nil(t, auth)
	require.Equal(t, sql.ErrNoRows, err)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestFatSecretAuthRepository_SaveAuth_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewFatSecretAuthRepository(sqlxDB)

	ctx := context.Background()
	auth := &models.FatSecretAuth{
		UserID:       1,
		AccessToken:  "access_token",
		AccessSecret: "access_secret",
	}

	mock.ExpectExec(`INSERT INTO FatsecretAuth`).
		WithArgs(auth.UserID, auth.AccessToken, auth.AccessSecret).
		WillReturnError(errors.New("db exec error"))

	err = repo.SaveAuth(ctx, auth)
	require.Error(t, err)
	require.Contains(t, err.Error(), "db exec error")

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}

func TestFatSecretAuthRepository_GetAuth_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewFatSecretAuthRepository(sqlxDB)

	ctx := context.Background()
	userID := 1

	mock.ExpectQuery(`SELECT user_id, access_token, access_secret FROM FatsecretAuth WHERE user_id = \$1`).
		WithArgs(userID).
		WillReturnError(sql.ErrNoRows)

	auth, err := repo.GetAuth(ctx, userID)
	require.Error(t, err)
	require.Nil(t, auth)
	require.Equal(t, sql.ErrNoRows, err)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}
