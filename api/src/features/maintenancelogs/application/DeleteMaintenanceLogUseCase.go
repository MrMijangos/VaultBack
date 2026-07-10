package application

import (
	"context"

	assetrepositories "vault/src/features/assets/domain/repositories"
	"vault/src/features/maintenancelogs/domain/repositories"
)

type DeleteMaintenanceLogUseCase struct {
	repo      repositories.MaintenanceLogRepository
	assetRepo assetrepositories.AssetRepository
}

func NewDeleteMaintenanceLogUseCase(repo repositories.MaintenanceLogRepository, assetRepo assetrepositories.AssetRepository) *DeleteMaintenanceLogUseCase {
	return &DeleteMaintenanceLogUseCase{repo: repo, assetRepo: assetRepo}
}

func (uc *DeleteMaintenanceLogUseCase) Execute(ctx context.Context, id string, userID string) error {
	existing, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	asset, err := uc.assetRepo.FindByID(ctx, existing.AssetID)
	if err != nil {
		return err
	}
	if asset.UserID != userID {
		return repositories.ErrMaintenanceLogNotFound
	}

	return uc.repo.Delete(ctx, id)
}
