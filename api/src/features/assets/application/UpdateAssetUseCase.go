package application

import (
	"context"
	"log"
	"time"

	"vault/src/core/eventbus"
	"vault/src/features/assets/domain/dto/request"
	"vault/src/features/assets/domain/dto/response"
	"vault/src/features/assets/domain/entities"
	"vault/src/features/assets/domain/repositories"
)

type UpdateAssetUseCase struct {
	repo      repositories.AssetRepository
	publisher eventbus.Publisher
}

func NewUpdateAssetUseCase(repo repositories.AssetRepository, publisher eventbus.Publisher) *UpdateAssetUseCase {
	return &UpdateAssetUseCase{repo: repo, publisher: publisher}
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
		Name:            req.Name,
		Category:        req.Category,
		Brand:           req.Brand,
		PurchaseValue:   req.PurchaseValue,
		Condition:       req.Condition,
		PurchaseDate:    purchaseDate,
		StoreOrigin:     req.StoreOrigin,
		Notes:           req.Notes,
		IsForSale:       req.IsForSale,
		SalePrice:       req.SalePrice,
		SaleDescription: req.SaleDescription,
		Size:            req.Size,
	}

	updated, err := uc.repo.Update(ctx, id, userID, asset)
	if err != nil {
		return response.AssetResponse{}, err
	}

	if err := uc.publisher.Publish(ctx, "asset.updated", updated.UserID, updated.ID, nil); err != nil {
		log.Printf("no se pudo publicar asset.updated: %v", err)
	}

	return response.FromEntity(updated, nil), nil
}
