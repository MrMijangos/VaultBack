package repositories

import (
	"context"
	"errors"

	"vault/src/features/assetcomments/domain/entities"
)

var ErrAssetCommentNotFound = errors.New("el comentario no existe")

type AssetCommentRepository interface {
	Create(ctx context.Context, comment entities.AssetComment) (entities.AssetComment, error)
	FindByAssetID(ctx context.Context, assetID string) ([]entities.AssetComment, error)
	Delete(ctx context.Context, id string, userID string) error
}
