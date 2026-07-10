package application

import (
	"context"
	"errors"

	"vault/src/core/security"
	"vault/src/features/auth/domain/dto/request"
	"vault/src/features/auth/domain/dto/response"
	"vault/src/features/auth/domain/repositories"
)

var ErrInvalidCredentials = errors.New("correo o contraseña incorrectos")

type LoginUseCase struct {
	repo      repositories.AuthRepository
	jwtSecret string
}

func NewLoginUseCase(repo repositories.AuthRepository, jwtSecret string) *LoginUseCase {
	return &LoginUseCase{repo: repo, jwtSecret: jwtSecret}
}

func (uc *LoginUseCase) Execute(ctx context.Context, req request.LoginRequest) (response.LoginResponse, string, error) {
	if err := req.Validate(); err != nil {
		return response.LoginResponse{}, "", err
	}

	creds, err := uc.repo.FindCredentialsByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, repositories.ErrCredentialsNotFound) {
			return response.LoginResponse{}, "", ErrInvalidCredentials
		}
		return response.LoginResponse{}, "", err
	}

	valid, err := security.VerifyPassword(creds.PasswordHash, req.Password)
	if err != nil || !valid {
		return response.LoginResponse{}, "", ErrInvalidCredentials
	}

	token, err := security.GenerateToken(creds.UserID, creds.Role, uc.jwtSecret)
	if err != nil {
		return response.LoginResponse{}, "", err
	}

	return response.FromCredentials(creds), token, nil
}
