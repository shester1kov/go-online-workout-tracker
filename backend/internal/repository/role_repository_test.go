package repository

import (
	"backend/internal/models"
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetRoleByID(t *testing.T) {

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	repo := NewRoleRepository(sqlxDB)

	roleID := 1
	expectedRole := &models.Role{
		ID:          roleID,
		Name:        "admin",
		Description: "Administrator role",
	}

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, name, description
		FROM Roles
		WHERE id = $1
		AND is_active = TRUE`,
	)).WithArgs(roleID).WillReturnRows(
		sqlmock.NewRows([]string{"id", "name", "description"}).
			AddRow(expectedRole.ID, expectedRole.Name, expectedRole.Description),
	)

	role, err := repo.GetRoleByID(context.Background(), roleID)
	assert.NoError(t, err)
	assert.Equal(t, expectedRole, role)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRoleRepository_GetRoleByID_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	repo := NewRoleRepository(sqlxDB)

	ctx := context.Background()
	roleID := 1

	mock.ExpectQuery(`SELECT id, name, description FROM Roles WHERE id = \$1 AND is_active = TRUE`).
		WithArgs(roleID).
		WillReturnError(errors.New("some db error"))

	role, err := repo.GetRoleByID(ctx, roleID)
	require.Error(t, err)
	require.Nil(t, role)

	err = mock.ExpectationsWereMet()
	require.NoError(t, err)
}
