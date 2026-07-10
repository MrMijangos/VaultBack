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

type UpdateMaintenanceLogUseCase struct {
	repo      repositories.MaintenanceLogRepository
	assetRepo assetrepositories.AssetRepository
}

func NewUpdateMaintenanceLogUseCase(repo repositories.MaintenanceLogRepository, assetRepo assetrepositories.AssetRepository) *UpdateMaintenanceLogUseCase {
	return &UpdateMaintenanceLogUseCase{repo: repo, assetRepo: assetRepo}
}

func (uc *UpdateMaintenanceLogUseCase) Execute(ctx context.Context, id string, userID string, req request.UpdateMaintenanceLogRequest) (response.MaintenanceLogResponse, error) {
	if err := req.Validate(); err != nil {
		return response.MaintenanceLogResponse{}, err
	}

	existing, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return response.MaintenanceLogResponse{}, err
	}

	asset, err := uc.assetRepo.FindByID(ctx, existing.AssetID)
	if err != nil {
		return response.MaintenanceLogResponse{}, err
	}
	if asset.UserID != userID {
		return response.MaintenanceLogResponse{}, repositories.ErrMaintenanceLogNotFound
	}

	var performedAt *time.Time
	if req.PerformedAt != "" {
		parsed, err := time.Parse("2006-01-02", req.PerformedAt)
		if err != nil {
			return response.MaintenanceLogResponse{}, err
		}
		performedAt = &parsed
	}

	updated, err := uc.repo.Update(ctx, id, entities.MaintenanceLog{
		Type:        req.Type,
		Subtype:     req.Subtype,
		Cost:        req.Cost,
		PerformedAt: performedAt,
		Notes:       req.Notes,
	})
	if err != nil {
		return response.MaintenanceLogResponse{}, err
	}

	return response.FromEntity(updated), nil
}
