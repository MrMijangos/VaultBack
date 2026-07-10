package application

import (
	"context"

	"vault/src/features/businesses/domain/dto/response"
	"vault/src/features/businesses/domain/repositories"
)

type GetBusinessByIdUseCase struct {
	repo repositories.BusinessRepository
}

func NewGetBusinessByIdUseCase(repo repositories.BusinessRepository) *GetBusinessByIdUseCase {
	return &GetBusinessByIdUseCase{repo: repo}
}

func (uc *GetBusinessByIdUseCase) Execute(ctx context.Context, id string) (response.BusinessResponse, error) {
	b, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return response.BusinessResponse{}, err
	}
	return response.FromEntity(b), nil
}
