package repository

import (
	"backend/internal/models"
	"context"
	"database/sql"
)

type RoleRepository struct {
	db *sql.DB
}

func NewRoleRepository(db *sql.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

func (r *RoleRepository) GetRoleByID(ctx context.Context, roleID int) (*models.Role, error) {
	query := `SELECT id, name, description
	FROM Roles
	WHERE id = $1
	AND is_active = TRUE`
	role := &models.Role{}

	err := r.db.QueryRowContext(ctx, query, roleID).Scan(
		&role.ID,
		&role.Name,
		&role.Description,
	)
	return role, err

}
