package auth

import (
	"context"
	"time"

	"payslip-generation-system/internal/common/errors"
	repo "payslip-generation-system/internal/repositories/user"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

//go:generate mockgen -source=service.go -package=mock -destination=mock/service_mock.go
type AuthServiceProvider interface {
	Login(ctx context.Context, username string, password string) (string, error)
}

type authService struct {
	repo repo.UserRepositoryProvider
	jwtSecret []byte
}

func NewAuthService(
	userRepo repo.UserRepositoryProvider,
	secret []byte,
) AuthServiceProvider {
	return &authService{
		repo: userRepo,
		jwtSecret: secret,
	}
}

func (s *authService) Login(ctx context.Context, username string, password string) (string, error) {
	user, err := s.repo.GetUserByUsername(ctx, username)
    if err != nil {
        return "", errors.New("invalid credentials")
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
        return "", errors.New("invalid credentials")
    }

    claims := jwt.MapClaims{
        "user_id":  user.ID,
        "is_admin": user.IsAdmin,
        "exp":      time.Now().Add(24 * time.Hour).Unix(),
        "iat":      time.Now().Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(s.jwtSecret)
}

