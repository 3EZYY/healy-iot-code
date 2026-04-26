package usecase

import (
	"context"
	"errors"

	"github.com/rafif/healy-backend/internal/domain"
	"github.com/rafif/healy-backend/internal/repository/interfaces"
	"github.com/rafif/healy-backend/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase interface {
	Login(ctx context.Context, req domain.LoginRequest) (domain.LoginResponse, error)
}

type authUsecase struct {
	userRepo       interfaces.UserRepository
	tokenGenerator jwt.TokenGenerator
}

func NewAuthUsecase(userRepo interfaces.UserRepository, tokenGenerator jwt.TokenGenerator) AuthUsecase {
	return &authUsecase{
		userRepo:       userRepo,
		tokenGenerator: tokenGenerator,
	}
}

func (u *authUsecase) Login(ctx context.Context, req domain.LoginRequest) (domain.LoginResponse, error) {
	// 1. Fetch user by username
	user, err := u.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		if err.Error() == "user not found" {
			return domain.LoginResponse{}, errors.New("invalid username or password")
		}
		return domain.LoginResponse{}, err
	}

	// 2. Compare passwords
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return domain.LoginResponse{}, errors.New("invalid username or password")
	}

	// 3. Generate JWT token
	tokenStr, expiresAt, err := u.tokenGenerator.GenerateToken(user)
	if err != nil {
		return domain.LoginResponse{}, err
	}

	return domain.LoginResponse{
		Token:     tokenStr,
		ExpiresAt: expiresAt,
	}, nil
}
