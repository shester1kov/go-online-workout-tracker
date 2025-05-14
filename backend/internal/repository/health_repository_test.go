package repository

import (
	"context"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func TestDBHealth_Check_Success(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("ошибка при создании sqlmock: %s", err)
	}
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewDBHealthRepository(sqlxDB)

	ctx := context.Background()

	mock.ExpectPing().WillReturnError(nil)

	err = repo.Check(ctx)
	if err != nil {
		t.Errorf("ожидалась успешная проверка, получили ошибку: %s", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("ожидания не выполнены: %s", err)
	}
}

func TestDBHealth_Check_Failure(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("ошибка при создании sqlmock: %s", err)
	}
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewDBHealthRepository(sqlxDB)

	ctx := context.Background()

	errPing := fmt.Errorf("ping failed")
	mock.ExpectPing().WillReturnError(errPing)

	err = repo.Check(ctx)
	if err == nil {
		t.Errorf("ожидалась ошибка, получили nil")
	}
	if err.Error() != errPing.Error() {
		t.Errorf("ожидалась ошибка '%s', получили '%s'", errPing, err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("ожидания не выполнены: %s", err)
	}
}
