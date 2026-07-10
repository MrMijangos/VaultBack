package application

import (
	"context"
	"time"

	"vault/src/features/assets/domain/dto/request"
	"vault/src/features/assets/domain/dto/response"
	"vault/src/features/assets/domain/entities"
	"vault/src/features/assets/domain/repositories"
)

type UpdateAssetUseCase struct {
	repo repositories.AssetRepository
}

func NewUpdateAssetUseCase(repo repositories.AssetRepository) *UpdateAssetUseCase {
	return &UpdateAssetUseCase{repo: repo}
}

func (uc *UpdateAssetUseCase) Execute(ctx context.Context, id string, userID string, req request.UpdateAssetRequest) (response.AssetResponse, error) {
	if err := req.Validate(); err != nil {
		return response.AssetResponse{}, err
	}

	var purchaseDate *time.Time
	if req.PurchaseDate != "" {
		parsed, err := time.Parse("2006-01-02", req.PurchaseDate)
		if err != nil {
			return response.AssetResponse{}, err
		}
		purchaseDate = &parsed
	}

	asset := entities.Asset{
		Name:          req.Name,
		Category:      req.Category,
		Brand:         req.Brand,
		PurchaseValue: req.PurchaseValue,
		Condition:     req.Condition,
		PurchaseDate:  purchaseDate,
		StoreOrigin:   req.StoreOrigin,
		Notes:         req.Notes,
	}

	updated, err := uc.repo.Update(ctx, id, userID, asset)
	if err != nil {
		return response.AssetResponse{}, err
	}

	return response.FromEntity(updated, nil), nil
}
