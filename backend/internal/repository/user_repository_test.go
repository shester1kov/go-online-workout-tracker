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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewUserRepository(sqlxDB)

	ctx := context.Background()
	user := &models.User{Username: "john", Email: "john@example.com", PasswordHash: "hash"}

	createdAt := time.Now()
	updatedAt := createdAt.Add(time.Hour)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO Users (username, email, password_hash)
	VALUES ($1, $2, $3)
	RETURNING id, created_at, updated_at`)).
		WithArgs(user.Username, user.Email, user.PasswordHash).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
			AddRow(1, createdAt, updatedAt))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id
	FROM Roles
	WHERE name = 'user'`)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2))

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO UserRoles (user_id, role_id)
	VALUES ($1, $2)`)).
		WithArgs(1, 2).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	err = repo.CreateUser(ctx, user)
	assert.NoError(t, err)
	assert.Equal(t, 1, user.ID)
	assert.WithinDuration(t, createdAt, user.CreatedAt, time.Second)
	assert.WithinDuration(t, updatedAt, user.UpdatedAt, time.Second)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAddUserRole(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewUserRepository(sqlxDB)

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO UserRoles (user_id, role_id)
	VALUES ($1, $2)`)).
		WithArgs(1, 3).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.AddUserRole(context.Background(), 1, 3)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewUserRepository(sqlxDB)

	email := "alice@example.com"
	existing := &models.User{ID: 5, Username: "alice", Email: email, PasswordHash: "pass", CreatedAt: time.Now()}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, username, email, password_hash, created_at
	FROM Users
	WHERE email = $1`)).
		WithArgs(email).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "password_hash", "created_at"}).
			AddRow(existing.ID, existing.Username, existing.Email, existing.PasswordHash, existing.CreatedAt))

	user, err := repo.GetUserByEmail(context.Background(), email)
	assert.NoError(t, err)
	assert.Equal(t, existing, user)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewUserRepository(sqlxDB)

	id := 7
	existing := &models.User{ID: id, Username: "bob", PasswordHash: "h", Email: "bob@example.com", CreatedAt: time.Now(), UpdatedAt: time.Now()}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, username, password_hash, email, created_at, updated_at
	FROM Users
	WHERE id = $1
	AND is_active = TRUE`)).
		WithArgs(id).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password_hash", "email", "created_at", "updated_at"}).
			AddRow(existing.ID, existing.Username, existing.PasswordHash, existing.Email, existing.CreatedAt, existing.UpdatedAt))

	user, err := repo.GetUserByID(context.Background(), id)
	assert.NoError(t, err)
	assert.Equal(t, existing, user)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserByUsername(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewUserRepository(sqlxDB)

	username := "charlie"
	existing := &models.User{ID: 9, Username: username, PasswordHash: "pwd", Email: "charlie@example.com", CreatedAt: time.Now(), UpdatedAt: time.Now()}

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, username, password_hash, email, created_at, updated_at
	FROM Users
	WHERE username = $1
	AND is_active = TRUE`)).
		WithArgs(username).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password_hash", "email", "created_at", "updated_at"}).
			AddRow(existing.ID, existing.Username, existing.PasswordHash, existing.Email, existing.CreatedAt, existing.UpdatedAt))

	user, err := repo.GetUserByUsername(context.Background(), username)
	assert.NoError(t, err)
	assert.Equal(t, existing, user)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserWithRolesByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewUserRepository(sqlxDB)

	id := 11
	userName := "dave"
	email := "dave@example.com"

	rows := sqlmock.NewRows([]string{"id", "username", "email", "id", "name", "description"}).
		AddRow(id, userName, email, 1, "admin", "desc1").
		AddRow(id, userName, email, 2, "user", "desc2")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT Users.id, Users.username, Users.email,
	Roles.id, Roles.name, Roles.description
	FROM Users
	INNER JOIN UserRoles AS ON Users.id = UserRoles.user_id
	INNER JOIN Roles AS ON UserRoles.role_id = Roles.id
	WHERE Users.id = $1`)).
		WithArgs(id).
		WillReturnRows(rows)

	user, err := repo.GetUserWithRolesByID(context.Background(), id)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, id, user.ID)
	assert.Equal(t, userName, user.Username)
	assert.Equal(t, email, user.Email)
	assert.Len(t, user.Roles, 2)
	assert.Equal(t, 1, user.Roles[0].ID)
	assert.Equal(t, "admin", user.Roles[0].Name)
	assert.Equal(t, 2, user.Roles[1].ID)
	assert.Equal(t, "user", user.Roles[1].Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserRoles(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewUserRepository(sqlxDB)

	id := 13

	rows := sqlmock.NewRows([]string{"id", "name", "description"}).
		AddRow(1, "admin", "desc1").
		AddRow(2, "user", "desc2")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT Roles.id, Roles.name, Roles.description
	FROM Roles
	INNER JOIN UserRoles ON Roles.id = UserRoles.role_id
	WHERE UserRoles.user_id = $1`)).
		WithArgs(id).
		WillReturnRows(rows)

	roles, err := repo.GetUserRoles(context.Background(), id)
	assert.NoError(t, err)
	assert.Len(t, *roles, 2)
	assert.Equal(t, 1, (*roles)[0].ID)
	assert.Equal(t, "admin", (*roles)[0].Name)
	assert.Equal(t, 2, (*roles)[1].ID)
	assert.Equal(t, "user", (*roles)[1].Name)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_CreateUser_ErrorBeginTx(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(sqlx.NewDb(db, "sqlmock"))

	mock.ExpectBegin().WillReturnError(errors.New("begin tx failed"))

	err = repo.CreateUser(context.Background(), &models.User{})
	assert.Error(t, err)
	assert.Equal(t, "begin tx failed", err.Error())
}

func TestUserRepository_CreateUser_ErrorInsertUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(sqlx.NewDb(db, "sqlmock"))

	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO Users").WillReturnError(errors.New("insert user failed"))
	mock.ExpectRollback()

	err = repo.CreateUser(context.Background(), &models.User{})
	assert.Error(t, err)
	assert.Equal(t, "insert user failed", err.Error())
}

func TestUserRepository_CreateUser_ErrorSelectRole(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(sqlx.NewDb(db, "sqlmock"))

	mock.ExpectBegin()
	// rows := sqlmock.NewRows([]string{"id"})
	mock.ExpectQuery("INSERT INTO Users").WillReturnRows(
		sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).AddRow(1, time.Now(), time.Now()),
	)
	mock.ExpectQuery("SELECT id FROM Roles").WillReturnError(errors.New("role query failed"))
	mock.ExpectRollback()

	err = repo.CreateUser(context.Background(), &models.User{})
	assert.Error(t, err)
	assert.Equal(t, "role query failed", err.Error())
}

func TestUserRepository_CreateUser_ErrorInsertUserRole(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(sqlx.NewDb(db, "sqlmock"))

	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO Users").WillReturnRows(
		sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).AddRow(1, time.Now(), time.Now()),
	)
	mock.ExpectQuery("SELECT id FROM Roles").WillReturnRows(
		sqlmock.NewRows([]string{"id"}).AddRow(2),
	)
	mock.ExpectExec("INSERT INTO UserRoles").WillReturnError(errors.New("insert role failed"))
	mock.ExpectRollback()

	err = repo.CreateUser(context.Background(), &models.User{})
	assert.Error(t, err)
	assert.Equal(t, "insert role failed", err.Error())
}

func TestUserRepository_CreateUser_ErrorCommit(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(sqlx.NewDb(db, "sqlmock"))

	mock.ExpectBegin()
	mock.ExpectQuery("INSERT INTO Users").WillReturnRows(
		sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).AddRow(1, time.Now(), time.Now()),
	)
	mock.ExpectQuery("SELECT id FROM Roles").WillReturnRows(
		sqlmock.NewRows([]string{"id"}).AddRow(2),
	)
	mock.ExpectExec("INSERT INTO UserRoles").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit().WillReturnError(errors.New("commit failed"))

	err = repo.CreateUser(context.Background(), &models.User{})
	assert.Error(t, err)
	assert.Equal(t, "commit failed", err.Error())
}

func TestUserRepository_AddUserRole_ErrorExec(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(sqlx.NewDb(db, "sqlmock"))

	mock.ExpectExec("INSERT INTO UserRoles").WillReturnError(errors.New("exec failed"))

	err = repo.AddUserRole(context.Background(), 1, 1)
	assert.Error(t, err)
	assert.Equal(t, "exec failed", err.Error())
}

func TestUserRepository_GetUserByEmail_ErrorQueryRow(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(sqlx.NewDb(db, "sqlmock"))

	mock.ExpectQuery("SELECT id, username, email, password_hash, created_at FROM Users WHERE email = \\$1").
		WithArgs("test@example.com").
		WillReturnError(errors.New("query failed"))

	user, err := repo.GetUserByEmail(context.Background(), "test@example.com")
	assert.Error(t, err)
	assert.Nil(t, user)
}

func TestUserRepository_GetUserByID_ErrorQueryRow(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(sqlx.NewDb(db, "sqlmock"))

	mock.ExpectQuery("SELECT id, username, password_hash, email, created_at, updated_at FROM Users WHERE id = \\$1 AND is_active = TRUE").
		WithArgs(1).
		WillReturnError(errors.New("query failed"))

	user, err := repo.GetUserByID(context.Background(), 1)
	assert.Error(t, err)
	assert.Nil(t, user)
}

func TestUserRepository_GetUserByUsername_ErrorQueryRow(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(sqlx.NewDb(db, "sqlmock"))

	mock.ExpectQuery("SELECT id, username, password_hash, email, created_at, updated_at FROM Users WHERE username = \\$1 AND is_active = TRUE").
		WithArgs("testuser").
		WillReturnError(errors.New("query failed"))

	user, err := repo.GetUserByUsername(context.Background(), "testuser")
	assert.Error(t, err)
	assert.Nil(t, user)
}

func TestUserRepository_GetUserWithRolesByID_ErrorQuery(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(sqlx.NewDb(db, "sqlmock"))

	mock.ExpectQuery("SELECT Users.id, Users.username, Users.email, Roles.id, Roles.name, Roles.description FROM Users INNER JOIN UserRoles").
		WithArgs(1).
		WillReturnError(errors.New("query failed"))

	user, err := repo.GetUserWithRolesByID(context.Background(), 1)
	assert.Error(t, err)
	assert.Nil(t, user)
}

func TestUserRepository_GetUserWithRolesByID_ErrorScan(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(sqlx.NewDb(db, "sqlmock"))

	rows := sqlmock.NewRows([]string{"Users.id", "Users.username", "Users.email", "Roles.id", "Roles.name", "Roles.description"}).
		AddRow(1, "user1", "email@example.com", 1, "role1", nil).
		AddRow(nil, nil, nil, nil, nil, nil) // заставим скан вернуть ошибку

	mock.ExpectQuery("SELECT Users.id, Users.username, Users.email, Roles.id, Roles.name, Roles.description FROM Users INNER JOIN UserRoles").
		WithArgs(1).
		WillReturnRows(rows)

	user, err := repo.GetUserWithRolesByID(context.Background(), 1)
	assert.Error(t, err)
	assert.Nil(t, user)
}

func TestUserRepository_GetUserRoles_ErrorQuery(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(sqlx.NewDb(db, "sqlmock"))

	mock.ExpectQuery("SELECT Roles.id, Roles.name, Roles.description FROM Roles INNER JOIN UserRoles").
		WithArgs(1).
		WillReturnError(errors.New("query failed"))

	roles, err := repo.GetUserRoles(context.Background(), 1)
	assert.Error(t, err)
	assert.Nil(t, roles)
}

func TestUserRepository_GetUserRoles_ErrorRowsScan(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewUserRepository(sqlx.NewDb(db, "sqlmock"))

	rows := sqlmock.NewRows([]string{"id", "name", "description"}).
		AddRow(1, "role1", "desc").
		AddRow(nil, nil, nil) // ошибочный row

	mock.ExpectQuery("SELECT Roles.id, Roles.name, Roles.description FROM Roles INNER JOIN UserRoles").
		WithArgs(1).
		WillReturnRows(rows)

	roles, err := repo.GetUserRoles(context.Background(), 1)
	assert.Error(t, err)
	assert.Nil(t, roles)
}
