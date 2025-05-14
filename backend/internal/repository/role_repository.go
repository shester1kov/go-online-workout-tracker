package repository

import (
	"backend/internal/models"
	"context"
	"log"

	"github.com/jmoiron/sqlx"
)

type RoleRepository struct {
	db *sqlx.DB
}

func NewRoleRepository(db *sqlx.DB) *RoleRepository {
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

	if err != nil {
		log.Println("Failed to get role by ID:", err)
		return nil, err
	}
	return role, err

}
