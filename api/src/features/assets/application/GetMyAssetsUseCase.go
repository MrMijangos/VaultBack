package application

import (
	"context"

	"vault/src/features/assets/domain/dto/response"
	"vault/src/features/assets/domain/repositories"
)

type GetMyAssetsUseCase struct {
	repo repositories.AssetRepository
}

func NewGetMyAssetsUseCase(repo repositories.AssetRepository) *GetMyAssetsUseCase {
	return &GetMyAssetsUseCase{repo: repo}
}

func (uc *GetMyAssetsUseCase) Execute(ctx context.Context, userID string) ([]response.AssetResponse, error) {
	assets, err := uc.repo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	list := make([]response.AssetResponse, 0, len(assets))
	for _, a := range assets {
		photos, err := uc.repo.FindPhotosByAssetID(ctx, a.ID)
		if err != nil {
			return nil, err
		}
		list = append(list, response.FromEntity(a, photos))
	}
	return list, nil
}
