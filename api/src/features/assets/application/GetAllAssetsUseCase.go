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

// Execute trae las fotos de cada activo por separado -- selectAssetsQuery
// no las incluye (photos vive en su propia tabla), así que sin esto la
// lista completa (usada por el Marketplace y "mis activos") siempre
// devolvía photos: [] aunque el activo ya tuviera fotos subidas.
func (uc *GetAllAssetsUseCase) Execute(ctx context.Context) ([]response.AssetResponse, error) {
	assets, err := uc.repo.FindAll(ctx)
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
