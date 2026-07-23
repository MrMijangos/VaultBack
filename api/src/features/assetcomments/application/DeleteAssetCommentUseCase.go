package application

import (
	"context"

	"vault/src/features/assetcomments/domain/repositories"
)

type DeleteAssetCommentUseCase struct {
	repo repositories.AssetCommentRepository
}

func NewDeleteAssetCommentUseCase(repo repositories.AssetCommentRepository) *DeleteAssetCommentUseCase {
	return &DeleteAssetCommentUseCase{repo: repo}
}

func (uc *DeleteAssetCommentUseCase) Execute(ctx context.Context, id string, userID string) error {
	return uc.repo.Delete(ctx, id, userID)
}
