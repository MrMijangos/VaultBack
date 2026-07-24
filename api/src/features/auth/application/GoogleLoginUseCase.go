package application

import (
	"context"

	"vault/src/core/security"
	"vault/src/features/auth/domain/dto/request"
	"vault/src/features/auth/domain/dto/response"
	"vault/src/features/auth/domain/repositories"
	"vault/src/features/auth/domain/services"
)

type GoogleLoginUseCase struct {
	repo      repositories.AuthRepository
	verifier  services.GoogleTokenVerifier
	jwtSecret string
}

func NewGoogleLoginUseCase(repo repositories.AuthRepository, verifier services.GoogleTokenVerifier, jwtSecret string) *GoogleLoginUseCase {
	return &GoogleLoginUseCase{repo: repo, verifier: verifier, jwtSecret: jwtSecret}
}

func (uc *GoogleLoginUseCase) Execute(ctx context.Context, req request.GoogleLoginRequest) (response.LoginResponse, string, error) {
	if err := req.Validate(); err != nil {
		return response.LoginResponse{}, "", err
	}

	profile, err := uc.verifier.Verify(ctx, req.IDToken)
	if err != nil {
		return response.LoginResponse{}, "", err
	}

	name := profile.Name
	if name == "" {
		name = profile.Email
	}

	creds, isNew, err := uc.repo.FindOrCreateByEmail(ctx, profile.Email, name, profile.Picture)
	if err != nil {
		return response.LoginResponse{}, "", err
	}

	token, err := security.GenerateToken(creds.UserID, creds.Role, uc.jwtSecret)
	if err != nil {
		return response.LoginResponse{}, "", err
	}

	return response.FromGoogleLogin(creds, isNew), token, nil
}
