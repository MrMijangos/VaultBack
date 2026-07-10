package application

import (
	"context"

	"vault/src/features/assets/domain/repositories"
)

type DeleteAssetUseCase struct {
	repo repositories.AssetRepository
}

func NewDeleteAssetUseCase(repo repositories.AssetRepository) *DeleteAssetUseCase {
	return &DeleteAssetUseCase{repo: repo}
}

func (uc *DeleteAssetUseCase) Execute(ctx context.Context, id string, userID string) error {
	return uc.repo.Delete(ctx, id, userID)
}
