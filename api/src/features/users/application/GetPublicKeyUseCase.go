package application

import (
	"context"

	"vault/src/features/users/domain/dto/response"
	"vault/src/features/users/domain/repositories"
)

type GetPublicKeyUseCase struct {
	repo repositories.UserRepository
}

func NewGetPublicKeyUseCase(repo repositories.UserRepository) *GetPublicKeyUseCase {
	return &GetPublicKeyUseCase{repo: repo}
}

func (uc *GetPublicKeyUseCase) Execute(ctx context.Context, userID string) (response.PublicKeyResponse, error) {
	publicKey, err := uc.repo.GetPublicKey(ctx, userID)
	if err != nil {
		return response.PublicKeyResponse{}, err
	}
	return response.PublicKeyResponse{UserID: userID, PublicKey: publicKey}, nil
}
