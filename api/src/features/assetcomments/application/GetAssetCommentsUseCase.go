package application

import (
	"context"

	"vault/src/features/assetcomments/domain/dto/response"
	"vault/src/features/assetcomments/domain/repositories"
)

type GetAssetCommentsUseCase struct {
	repo repositories.AssetCommentRepository
}

func NewGetAssetCommentsUseCase(repo repositories.AssetCommentRepository) *GetAssetCommentsUseCase {
	return &GetAssetCommentsUseCase{repo: repo}
}

func (uc *GetAssetCommentsUseCase) Execute(ctx context.Context, assetID string) ([]response.AssetCommentResponse, error) {
	list, err := uc.repo.FindByAssetID(ctx, assetID)
	if err != nil {
		return nil, err
	}
	return response.FromEntities(list), nil
}
