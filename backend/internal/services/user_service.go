package services

import (
	"backend/internal/apperrors"
	"backend/internal/models"
	"backend/internal/repository"
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/lib/pq"
)

type UserService struct {
	userRepo *repository.UserRepository
	roleRepo *repository.RoleRepository
}

func NewUserService(userRepo *repository.UserRepository, roleRepo *repository.RoleRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
		roleRepo: roleRepo,
	}
}

func (s *UserService) GetUserByID(ctx context.Context, userID int) (*models.User, error) {

	user, err := s.userRepo.GetUserByID(ctx, userID)
	if user == nil {
		return nil, &apperrors.AppError{
			Code:    http.StatusNotFound,
			Message: "User not found",
		}
	}

	if err != nil {
		switch {
		case errors.Is(err, context.DeadlineExceeded):
			return nil, &apperrors.AppError{
				Code:    http.StatusGatewayTimeout,
				Message: "Request timeout",
			}
		case errors.Is(err, context.Canceled):
			return nil, &apperrors.AppError{
				Code:    http.StatusBadRequest,
				Message: "Request cancelled",
			}
		case errors.Is(err, sql.ErrNoRows):
			return nil, &apperrors.AppError{
				Code:    http.StatusNotFound,
				Message: "User not found",
			}
		default:
			return nil, &apperrors.AppError{
				Code:    http.StatusInternalServerError,
				Message: "Failed to get user",
			}
		}
	}
	return user, nil
}

func (s *UserService) AddRoleToUser(ctx context.Context, id int, req *models.AddRoleToUserRequest) (*models.User, error) {

	if err := s.userRepo.AddUserRole(ctx, id, req.RoleID); err != nil {
		var pgErr *pq.Error
		switch {
		case errors.As(err, &pgErr) && pgErr.Code == apperrors.PgErrUniqueViolation:
			return nil, &apperrors.AppError{
				Code:    http.StatusConflict,
				Message: "User already has this role",
			}
		case errors.As(err, &pgErr) && pgErr.Code == apperrors.PgErrForeignKeyViolation:
			if strings.Contains(pgErr.Constraint, "user_id") {
				return nil, &apperrors.AppError{
					Code:    http.StatusNotFound,
					Message: "User not found",
				}
			} else if strings.Contains(pgErr.Constraint, "role_id") {
				return nil, &apperrors.AppError{
					Code:    http.StatusNotFound,
					Message: "Role not found",
				}
			}
		case errors.Is(err, context.DeadlineExceeded):
			return nil, &apperrors.AppError{
				Code:    http.StatusGatewayTimeout,
				Message: "Request timeout",
			}
		case errors.Is(err, context.Canceled):
			return nil, &apperrors.AppError{
				Code:    http.StatusBadRequest,
				Message: "Request cancelled",
			}
		default:
			return nil, &apperrors.AppError{
				Code:    http.StatusInternalServerError,
				Message: "Failed to add role",
			}
		}
	}

	return s.userRepo.GetUserByID(ctx, id)
}

func (s *UserService) GetUserRoles(ctx context.Context, userID int) (*[]models.Role, error) {
	roles, err := s.userRepo.GetUserRoles(ctx, userID)
	if err != nil {
		log.Println("Failed to get user roles:", err)
		return nil, err
	}
	return roles, nil
}
