package application

import (
	"context"

	"vault/src/features/businessservices/domain/dto/response"
	"vault/src/features/businessservices/domain/repositories"
)

type ListBusinessServicesUseCase struct {
	repo repositories.BusinessServiceRepository
}

func NewListBusinessServicesUseCase(repo repositories.BusinessServiceRepository) *ListBusinessServicesUseCase {
	return &ListBusinessServicesUseCase{repo: repo}
}

func (uc *ListBusinessServicesUseCase) Execute(ctx context.Context, businessID string) ([]response.BusinessServiceResponse, error) {
	list, err := uc.repo.ListByBusinessID(ctx, businessID)
	if err != nil {
		return nil, err
	}
	return response.FromEntities(list), nil
}
