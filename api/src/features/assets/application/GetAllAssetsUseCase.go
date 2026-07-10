package application

import (
	"context"

	"vault/src/features/assets/domain/dto/response"
	"vault/src/features/assets/domain/repositories"
)

type GetAllAssetsUseCase struct {
	repo repositories.AssetRepository
}

func NewGetAllAssetsUseCase(repo repositories.AssetRepository) *GetAllAssetsUseCase {
	return &GetAllAssetsUseCase{repo: repo}
}

func (uc *GetAllAssetsUseCase) Execute(ctx context.Context) ([]response.AssetResponse, error) {
	assets, err := uc.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return response.FromEntities(assets), nil
}
