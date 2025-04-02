package services

import (
	"backend/internal/apperrors"
	"backend/internal/auth"
	"backend/internal/models"
	"backend/internal/repository"
	"backend/internal/utils"
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/lib/pq"
)

type AuthService struct {
	userRepo   *repository.UserRepository
	jwtManager *auth.JWTManager
}

func NewAuthService(userRepo *repository.UserRepository, jwtManager *auth.JWTManager) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		jwtManager: jwtManager,
	}
}

func (s *AuthService) Register(ctx context.Context, req *models.UserRegisterRequest) (*models.User, error) {

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		log.Println("Failed to hash password:", err)
		return nil, &apperrors.AppError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to hash password",
		}
	}

	user := &models.User{
		Username:     req.Username,
		PasswordHash: hashedPassword,
		Email:        req.Email,
	}

	if err := s.userRepo.CreateUser(ctx, user); err != nil {
		log.Println("Failed to create user:", err)
		var pgErr *pq.Error
		switch {
		case errors.As(err, &pgErr) && pgErr.Code == apperrors.PgErrUniqueViolation:
			constraint := pgErr.Constraint
			if strings.Contains(constraint, "username") {
				return nil, &apperrors.AppError{
					Code:    http.StatusConflict,
					Message: "Username already exists",
				}
			} else if strings.Contains(constraint, "email") {
				return nil, &apperrors.AppError{
					Code:    http.StatusConflict,
					Message: "Email already exists",
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
		case errors.Is(err, sql.ErrNoRows):
			return nil, &apperrors.AppError{
				Code:    http.StatusInternalServerError,
				Message: "Default role 'user' not found in database",
			}
		default:
			log.Println("Unhandled error:", err)
			return nil, &apperrors.AppError{
				Code:    http.StatusInternalServerError,
				Message: "Failed to create user",
			}
		}
	}

	return user, nil
}

func (s *AuthService) Login(ctx context.Context, req *models.UserAuthRequest) (string, *models.User, error) {

	user, err := s.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return "", nil, &apperrors.AppError{
				Code:    http.StatusUnauthorized,
				Message: "Invalid email or password",
			}
		case errors.Is(err, context.DeadlineExceeded):
			return "", nil, &apperrors.AppError{
				Code:    http.StatusGatewayTimeout,
				Message: "Request timeout",
			}
		case errors.Is(err, context.Canceled):
			return "", nil, &apperrors.AppError{
				Code:    http.StatusBadRequest,
				Message: "Request cancelled",
			}
		default:
			log.Println("Unhandled error:", err)
			return "", nil, &apperrors.AppError{
				Code:    http.StatusInternalServerError,
				Message: "Internal server error",
			}
		}
	}

	if err := utils.CheckPassword(req.Password, user.PasswordHash); err != nil {
		return "", nil, &apperrors.AppError{
			Code:    http.StatusUnauthorized,
			Message: "Invalid email or password",
		}
	}

	token, err := s.jwtManager.Generate(user)
	if err != nil {
		return "", nil, &apperrors.AppError{
			Code:    http.StatusInternalServerError,
			Message: "Failed to generate token",
		}
	}

	return token, user, nil
}
