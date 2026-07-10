package application

import (
	"context"
	"time"

	"vault/src/features/assets/domain/dto/request"
	"vault/src/features/assets/domain/dto/response"
	"vault/src/features/assets/domain/entities"
	"vault/src/features/assets/domain/repositories"
)

type CreateAssetUseCase struct {
	repo repositories.AssetRepository
}

func NewCreateAssetUseCase(repo repositories.AssetRepository) *CreateAssetUseCase {
	return &CreateAssetUseCase{repo: repo}
}

func (uc *CreateAssetUseCase) Execute(ctx context.Context, userID string, req request.CreateAssetRequest) (response.AssetResponse, error) {
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
		UserID:        userID,
		Name:          req.Name,
		Category:      req.Category,
		Brand:         req.Brand,
		PurchaseValue: req.PurchaseValue,
		Condition:     req.Condition,
		PurchaseDate:  purchaseDate,
		StoreOrigin:   req.StoreOrigin,
		Notes:         req.Notes,
	}

	created, err := uc.repo.Create(ctx, asset)
	if err != nil {
		return response.AssetResponse{}, err
	}

	return response.FromEntity(created, nil), nil
}
