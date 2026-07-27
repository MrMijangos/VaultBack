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

type CreateAssetUseCase struct {
	repo      repositories.AssetRepository
	publisher eventbus.Publisher
}

func NewCreateAssetUseCase(repo repositories.AssetRepository, publisher eventbus.Publisher) *CreateAssetUseCase {
	return &CreateAssetUseCase{repo: repo, publisher: publisher}
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
		UserID:          userID,
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

	created, err := uc.repo.Create(ctx, asset)
	if err != nil {
		return response.AssetResponse{}, err
	}

	// asset.updated: vault-ai-service usa la misma routing key para
	// cualquier cambio que afecte el perfil ML del usuario.
	if err := uc.publisher.Publish(ctx, "asset.updated", created.UserID, created.ID, nil); err != nil {
		log.Printf("no se pudo publicar asset.updated: %v", err)
	}

	// asset.created: lo consume vault-blockchain para certificar el activo
	// en Vara Network -- necesita name/category/created_at para calcular el
	// hash del activo, por eso no reusa el shape de Publish/asset.updated.
	assetCreatedPayload := struct {
		EventType string    `json:"event_type"`
		AssetID   string    `json:"asset_id"`
		UserID    string    `json:"user_id"`
		Name      string    `json:"name"`
		Category  string    `json:"category"`
		CreatedAt time.Time `json:"created_at"`
	}{
		EventType: "asset.created",
		AssetID:   created.ID,
		UserID:    created.UserID,
		Name:      created.Name,
		Category:  created.Category,
		CreatedAt: created.CreatedAt,
	}
	if err := uc.publisher.PublishEvent(ctx, "asset.created", assetCreatedPayload); err != nil {
		log.Printf("no se pudo publicar asset.created: %v", err)
	}

	return response.FromEntity(created, nil), nil
}
