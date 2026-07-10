package repositories

import (
	"context"
	"errors"

	"vault/src/features/assets/domain/entities"
)

var ErrAssetNotFound = errors.New("el producto no existe")

type AssetRepository interface {
	Create(ctx context.Context, asset entities.Asset) (entities.Asset, error)
	FindAll(ctx context.Context) ([]entities.Asset, error)
	FindByID(ctx context.Context, id string) (entities.Asset, error)
	Update(ctx context.Context, id string, userID string, asset entities.Asset) (entities.Asset, error)
	Delete(ctx context.Context, id string, userID string) error
	FindPhotosByAssetID(ctx context.Context, assetID string) ([]entities.AssetPhoto, error)
	AddPhoto(ctx context.Context, assetID string, url string) (entities.AssetPhoto, error)
}
