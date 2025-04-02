package repository

import (
	"context"
	"database/sql"
)

type DBHeathRepository struct {
	db *sql.DB
}

func NewDBHealthRepository(db *sql.DB) *DBHeathRepository {
	return &DBHeathRepository{db: db}
}

func (r *DBHeathRepository) Check(ctx context.Context) error {
	if err := r.db.PingContext(ctx); err != nil {
		return err
	}
	return nil
}
