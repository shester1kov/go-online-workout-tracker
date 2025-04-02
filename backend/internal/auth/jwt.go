package auth

import (
	"backend/internal/config"
	"backend/internal/models"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
	secretKey string
}

func InitJWTManager(envs *config.Envs) *JWTManager {
	if len(envs.JWTSecureKey) < 32 {
		panic("JWT secret key must be at least 32 bytes")
	}
	return &JWTManager{secretKey: envs.JWTSecureKey}
}

func (m *JWTManager) Generate(user *models.User) (string, error) {
	claims := &models.Claims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.secretKey))
}

func (m *JWTManager) Verify(tokenString string) (*models.Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&models.Claims{},
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(m.secretKey), nil
		},
	)
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*models.Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
