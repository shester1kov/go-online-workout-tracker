package repository

import (
	"backend/internal/models"
	"context"
	"log"

	"github.com/jmoiron/sqlx"
)

type FatSecretAuthRepository struct {
	db *sqlx.DB
}

func NewFatSecretAuthRepository(db *sqlx.DB) *FatSecretAuthRepository {
	return &FatSecretAuthRepository{db: db}
}

func (r *FatSecretAuthRepository) SaveTempAuth(ctx context.Context, auth *models.TempAuth) error {
	query := `INSERT INTO TempFatsecretAuth (user_id, request_token, request_secret)
         VALUES ($1, $2, $3)
         ON CONFLICT (user_id) DO UPDATE
         SET request_token = $2, request_secret = $3, created_at = NOW()`

	_, err := r.db.ExecContext(
		ctx,
		query,
		auth.UserID,
		auth.RequestToken,
		auth.RequestSecret,
	)

	if err != nil {
		log.Println("Save temporary auth error:", err)
		return err
	}

	return nil
}

func (r *FatSecretAuthRepository) GetTempAuth(ctx context.Context, token string) (*models.TempAuth, error) {
	var auth models.TempAuth

	query := `SELECT user_id, request_token, request_secret
	FROM TempFatsecretAuth
    WHERE request_token = $1`

	err := r.db.QueryRowContext(
		ctx,
		query,
		token,
	).Scan(
		&auth.UserID,
		&auth.RequestToken,
		&auth.RequestSecret,
	)
	if err != nil {
		log.Println("Get temporary auth error:", err)
		return nil, err
	}
	return &auth, nil
}

func (r *FatSecretAuthRepository) SaveAuth(ctx context.Context, auth *models.FatSecretAuth) error {
	query := `INSERT INTO FatsecretAuth (user_id, access_token, access_secret)
    VALUES ($1, $2, $3)
    ON CONFLICT (user_id) DO UPDATE
    SET access_token = $2, access_secret = $3, updated_at = NOW()`

	_, err := r.db.ExecContext(
		ctx,
		query,
		auth.UserID,
		auth.AccessToken,
		auth.AccessSecret,
	)

	if err != nil {
		log.Println("Save auth error:", err)
		return err
	}

	return nil
}

func (r *FatSecretAuthRepository) GetAuth(ctx context.Context, userID int) (*models.FatSecretAuth, error) {
	var auth models.FatSecretAuth

	query := `SELECT user_id, access_token, access_secret
	FROM FatsecretAuth
	WHERE user_id = $1`

	err := r.db.QueryRowContext(ctx,
		query,
		userID,
	).Scan(
		&auth.UserID,
		&auth.AccessToken,
		&auth.AccessSecret,
	)
	if err != nil {
		log.Println("Get auth error:", err)
		return nil, err
	}
	return &auth, nil
}
