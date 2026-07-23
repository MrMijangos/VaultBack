package application

import (
	"context"

	"vault/src/features/users/domain/dto/request"
	"vault/src/features/users/domain/repositories"
)

type SetPublicKeyUseCase struct {
	repo repositories.UserRepository
}

func NewSetPublicKeyUseCase(repo repositories.UserRepository) *SetPublicKeyUseCase {
	return &SetPublicKeyUseCase{repo: repo}
}

func (uc *SetPublicKeyUseCase) Execute(ctx context.Context, userID string, req request.SetPublicKeyRequest) error {
	if err := req.Validate(); err != nil {
		return err
	}
	return uc.repo.SetPublicKey(ctx, userID, req.PublicKey)
}
