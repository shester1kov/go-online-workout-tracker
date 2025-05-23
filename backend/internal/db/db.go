package db

import (
	"backend/internal/config"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewConnection(envs *config.Envs) (*sqlx.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		envs.PostgresHost,
		envs.PostgresPort,
		envs.PostgresUser,
		envs.PostgresPassword,
		envs.PostgresDb,
		envs.PostgresSslmode,
	)

	//fmt.Println(dsn)

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("database is unreachable: %w", err)
	}

	return db, nil
}
