package application

import (
	"context"
	"log"
	"time"

	"vault/src/core/eventbus"
	assetrepositories "vault/src/features/assets/domain/repositories"
	"vault/src/features/maintenancelogs/domain/dto/request"
	"vault/src/features/maintenancelogs/domain/dto/response"
	"vault/src/features/maintenancelogs/domain/entities"
	"vault/src/features/maintenancelogs/domain/repositories"
)

type CreateMaintenanceLogUseCase struct {
	repo      repositories.MaintenanceLogRepository
	assetRepo assetrepositories.AssetRepository
	publisher eventbus.Publisher
}

func NewCreateMaintenanceLogUseCase(repo repositories.MaintenanceLogRepository, assetRepo assetrepositories.AssetRepository, publisher eventbus.Publisher) *CreateMaintenanceLogUseCase {
	return &CreateMaintenanceLogUseCase{repo: repo, assetRepo: assetRepo, publisher: publisher}
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

	// maintenance.registered: lo consume vault-blockchain para certificar el
	// mantenimiento/restauración en Vara Network.
	maintenanceRegisteredPayload := struct {
		EventType     string `json:"event_type"`
		MaintenanceID string `json:"maintenance_id"`
		AssetID       string `json:"asset_id"`
		UserID        string `json:"user_id"`
		Type          string `json:"type"`
	}{
		EventType:     "maintenance.registered",
		MaintenanceID: created.ID,
		AssetID:       created.AssetID,
		UserID:        userID,
		Type:          created.Type,
	}
	if err := uc.publisher.PublishEvent(ctx, "maintenance.registered", maintenanceRegisteredPayload); err != nil {
		log.Printf("no se pudo publicar maintenance.registered: %v", err)
	}

	return response.FromEntity(created), nil
}
