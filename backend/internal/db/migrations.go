package db

import (
	"backend/internal/config"
	"fmt"
	"log"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func RunMigrations(db *sqlx.DB, envs *config.Envs) error {
	migrationsPath, err := filepath.Abs("migrations")
	if err != nil {
		return fmt.Errorf("abs path get error: %v", err)
	}

	migrationsPath = filepath.ToSlash(migrationsPath)

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
		envs.PostgresDb,
		driver,
	)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Println("migration failed", err)
		return fmt.Errorf("migration failed: %v", err)
	}

	log.Println("migrations applied successfully")
	return nil
}
