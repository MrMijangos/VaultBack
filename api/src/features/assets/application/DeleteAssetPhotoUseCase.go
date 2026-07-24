package application

import (
	"context"

	"vault/src/features/assets/domain/dto/response"
	"vault/src/features/assets/domain/repositories"
)

type DeleteAssetPhotoUseCase struct {
	repo repositories.AssetRepository
}

func NewDeleteAssetPhotoUseCase(repo repositories.AssetRepository) *DeleteAssetPhotoUseCase {
	return &DeleteAssetPhotoUseCase{repo: repo}
}

func (uc *DeleteAssetPhotoUseCase) Execute(ctx context.Context, assetID string, photoID string, userID string) (response.AssetResponse, error) {
	asset, err := uc.repo.FindByID(ctx, assetID)
	if err != nil {
		return response.AssetResponse{}, err
	}
	if asset.UserID != userID {
		return response.AssetResponse{}, repositories.ErrAssetNotFound
	}

	if err := uc.repo.DeletePhoto(ctx, photoID, assetID); err != nil {
		return response.AssetResponse{}, err
	}

	photos, err := uc.repo.FindPhotosByAssetID(ctx, assetID)
	if err != nil {
		return response.AssetResponse{}, err
	}

	return response.FromEntity(asset, photos), nil
}
