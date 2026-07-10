package application

import (
	"context"

	"vault/src/features/maintenancelogs/domain/dto/response"
	"vault/src/features/maintenancelogs/domain/repositories"
)

type GetLogsByAssetUseCase struct {
	repo repositories.MaintenanceLogRepository
}

func NewGetLogsByAssetUseCase(repo repositories.MaintenanceLogRepository) *GetLogsByAssetUseCase {
	return &GetLogsByAssetUseCase{repo: repo}
}

func (uc *GetLogsByAssetUseCase) Execute(ctx context.Context, assetID string) ([]response.MaintenanceLogResponse, error) {
	list, err := uc.repo.FindByAssetID(ctx, assetID)
	if err != nil {
		return nil, err
	}
	return response.FromEntities(list), nil
}
