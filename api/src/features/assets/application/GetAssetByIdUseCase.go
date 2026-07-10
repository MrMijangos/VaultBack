package application

import (
	"context"

	"vault/src/features/assets/domain/dto/response"
	"vault/src/features/assets/domain/repositories"
)

type GetAssetByIdUseCase struct {
	repo repositories.AssetRepository
}

func NewGetAssetByIdUseCase(repo repositories.AssetRepository) *GetAssetByIdUseCase {
	return &GetAssetByIdUseCase{repo: repo}
}

func (uc *GetAssetByIdUseCase) Execute(ctx context.Context, id string) (response.AssetResponse, error) {
	asset, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return response.AssetResponse{}, err
	}

	photos, err := uc.repo.FindPhotosByAssetID(ctx, id)
	if err != nil {
		return response.AssetResponse{}, err
	}

	return response.FromEntity(asset, photos), nil
}
