package application

import (
	"context"

	"vault/src/features/businesses/domain/dto/response"
	"vault/src/features/businesses/domain/repositories"
)

type GetAllBusinessesUseCase struct {
	repo repositories.BusinessRepository
}

func NewGetAllBusinessesUseCase(repo repositories.BusinessRepository) *GetAllBusinessesUseCase {
	return &GetAllBusinessesUseCase{repo: repo}
}

func (uc *GetAllBusinessesUseCase) Execute(ctx context.Context) ([]response.BusinessResponse, error) {
	list, err := uc.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return response.FromEntities(list), nil
}
