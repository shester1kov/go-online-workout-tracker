package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type DBHeathRepository struct {
	db *sqlx.DB
}

func NewDBHealthRepository(db *sqlx.DB) *DBHeathRepository {
	return &DBHeathRepository{db: db}
}

func (r *DBHeathRepository) Check(ctx context.Context) error {
	if err := r.db.PingContext(ctx); err != nil {
		return err
	}
	return nil
}
