package repository

import (
	"backend/internal/models"
	"context"
	"log"

	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		log.Println("Failed to begin transaction:", err)
		return err
	}

	defer func() {
		if err != nil {
			log.Println("Error, rollback transaction")
			tx.Rollback()
		}
	}()

	userQuery := `INSERT INTO Users (username, email, password_hash)
	VALUES ($1, $2, $3)
	RETURNING id, created_at, updated_at`

	err = tx.QueryRowContext(
		ctx,
		userQuery,
		user.Username,
		user.Email,
		user.PasswordHash,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		log.Println("Failed to create user:", err)
		return err
	}

	var roleID int
	roleQuery := `SELECT id
	FROM Roles
	WHERE name = 'user'`

	err = tx.QueryRowContext(ctx, roleQuery).Scan(&roleID)
	if err != nil {
		log.Println("Failed to get role:", err)
		return err
	}

	userRoleQuery := `INSERT INTO UserRoles (user_id, role_id)
	VALUES ($1, $2)`
	_, err = tx.Exec(userRoleQuery, user.ID, roleID)
	if err != nil {
		log.Println("Failed to add role:", err)
		return err
	}

	if err = tx.Commit(); err != nil {
		log.Println("Failed to commit transaction:", err)
		return err
	}

	return nil
}

func (r *UserRepository) AddUserRole(ctx context.Context, userID, roleID int) error {
	query := `INSERT INTO UserRoles (user_id, role_id)
	VALUES ($1, $2)`
	_, err := r.db.ExecContext(
		ctx,
		query,
		userID,
		roleID,
	)
	if err != nil {
		log.Println("Error add user role:", err)
		return err
	}
	return nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `SELECT id, username, email, password_hash, created_at
	FROM Users
	WHERE email = $1`
	user := &models.User{}

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
	)
	if err != nil {
		log.Println("Failed to get user by email:", err)
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, userID int) (*models.User, error) {
	query := `SELECT id, username, password_hash, email, created_at, updated_at
	FROM Users
	WHERE id = $1
	AND is_active = TRUE`
	user := &models.User{}

	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		log.Println("Failed to get user by id:", err)
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {

	query := `SELECT id, username, password_hash, email, created_at, updated_at
	FROM Users
	WHERE username = $1
	AND is_active = TRUE`
	user := &models.User{}

	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		log.Println("Failed to get user by username:", err)
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetUserWithRolesByID(ctx context.Context, userID int) (*models.User, error) {
	query := `SELECT Users.id, Users.username, Users.email,
	Roles.id, Roles.name, Roles.description
	FROM Users
	INNER JOIN UserRoles AS ON Users.id = UserRoles.user_id
	INNER JOIN Roles AS ON UserRoles.role_id = Roles.id
	WHERE Users.id = $1`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var user models.User
	var roles []models.Role

	for rows.Next() {
		var role models.Role
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&role.ID,
			&role.Name,
			&role.Description,
		)
		if err != nil {
			return nil, err
		}
		if role.ID != 0 {
			roles = append(roles, role)
		}
	}

	if user.ID == 0 {
		return nil, nil
	}

	user.Roles = roles
	return &user, nil
}

func (r *UserRepository) GetUserRoles(ctx context.Context, userID int) (*[]models.Role, error) {
	query := `SELECT Roles.id, Roles.name, Roles.description
	FROM Roles
	INNER JOIN UserRoles ON Roles.id = UserRoles.role_id
	WHERE UserRoles.user_id = $1`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		log.Println("Failed to get user roles:", err)
		return nil, err
	}
	defer rows.Close()

	var roles []models.Role

	for rows.Next() {
		var role models.Role
		err := rows.Scan(
			&role.ID,
			&role.Name,
			&role.Description,
		)
		if err != nil {
			log.Println("Failed to scan user role:", err)
			return nil, err
		}

		roles = append(roles, role)
	}

	if rows.Err() != nil {
		log.Println("Rows iteration error:", err)
		return nil, err
	}

	return &roles, nil
}
