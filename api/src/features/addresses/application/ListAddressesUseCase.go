package application

import (
	"context"

	"vault/src/features/addresses/domain/dto/response"
	"vault/src/features/addresses/domain/repositories"
)

type ListAddressesUseCase struct {
	repo repositories.AddressRepository
}

func NewListAddressesUseCase(repo repositories.AddressRepository) *ListAddressesUseCase {
	return &ListAddressesUseCase{repo: repo}
}

func (uc *ListAddressesUseCase) Execute(ctx context.Context, userID string) ([]response.AddressResponse, error) {
	list, err := uc.repo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return response.FromEntities(list), nil
}
