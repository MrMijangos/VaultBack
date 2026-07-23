package application

import (
	"context"

	"vault/src/features/addresses/domain/repositories"
)

type DeleteAddressUseCase struct {
	repo repositories.AddressRepository
}

func NewDeleteAddressUseCase(repo repositories.AddressRepository) *DeleteAddressUseCase {
	return &DeleteAddressUseCase{repo: repo}
}

func (uc *DeleteAddressUseCase) Execute(ctx context.Context, id string, userID string) error {
	return uc.repo.Delete(ctx, id, userID)
}
