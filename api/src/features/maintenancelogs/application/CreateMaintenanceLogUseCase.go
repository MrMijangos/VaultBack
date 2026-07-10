package application

import (
	"context"
	"time"

	assetrepositories "vault/src/features/assets/domain/repositories"
	"vault/src/features/maintenancelogs/domain/dto/request"
	"vault/src/features/maintenancelogs/domain/dto/response"
	"vault/src/features/maintenancelogs/domain/entities"
	"vault/src/features/maintenancelogs/domain/repositories"
)

type CreateMaintenanceLogUseCase struct {
	repo      repositories.MaintenanceLogRepository
	assetRepo assetrepositories.AssetRepository
}

func NewCreateMaintenanceLogUseCase(repo repositories.MaintenanceLogRepository, assetRepo assetrepositories.AssetRepository) *CreateMaintenanceLogUseCase {
	return &CreateMaintenanceLogUseCase{repo: repo, assetRepo: assetRepo}
}

func (uc *CreateMaintenanceLogUseCase) Execute(ctx context.Context, userID string, req request.CreateMaintenanceLogRequest) (response.MaintenanceLogResponse, error) {
	if err := req.Validate(); err != nil {
		return response.MaintenanceLogResponse{}, err
	}

	asset, err := uc.assetRepo.FindByID(ctx, req.AssetID)
	if err != nil {
		return response.MaintenanceLogResponse{}, err
	}
	if asset.UserID != userID {
		return response.MaintenanceLogResponse{}, assetrepositories.ErrAssetNotFound
	}

	var providerID *string
	if req.ProviderID != "" {
		providerID = &req.ProviderID
	}

	var performedAt *time.Time
	if req.PerformedAt != "" {
		parsed, err := time.Parse("2006-01-02", req.PerformedAt)
		if err != nil {
			return response.MaintenanceLogResponse{}, err
		}
		performedAt = &parsed
	}

	created, err := uc.repo.Create(ctx, entities.MaintenanceLog{
		AssetID:     req.AssetID,
		ProviderID:  providerID,
		Type:        req.Type,
		Subtype:     req.Subtype,
		Cost:        req.Cost,
		PerformedAt: performedAt,
		Notes:       req.Notes,
	})
	if err != nil {
		return response.MaintenanceLogResponse{}, err
	}

	return response.FromEntity(created), nil
}
