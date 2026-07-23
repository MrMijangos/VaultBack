package application

import (
	"context"

	"vault/src/features/addresses/domain/dto/response"
	"vault/src/features/addresses/domain/repositories"
)

type SetDefaultAddressUseCase struct {
	repo repositories.AddressRepository
}

func NewSetDefaultAddressUseCase(repo repositories.AddressRepository) *SetDefaultAddressUseCase {
	return &SetDefaultAddressUseCase{repo: repo}
}

func (uc *SetDefaultAddressUseCase) Execute(ctx context.Context, id string, userID string) (response.AddressResponse, error) {
	updated, err := uc.repo.SetDefault(ctx, id, userID)
	if err != nil {
		return response.AddressResponse{}, err
	}
	return response.FromEntity(updated), nil
}
