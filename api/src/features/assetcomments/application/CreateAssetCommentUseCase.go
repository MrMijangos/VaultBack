package application

import (
	"context"

	"github.com/google/uuid"

	"vault/src/core/moderation"
	"vault/src/features/assetcomments/domain/dto/request"
	"vault/src/features/assetcomments/domain/dto/response"
	"vault/src/features/assetcomments/domain/entities"
	"vault/src/features/assetcomments/domain/repositories"
)

type CreateAssetCommentUseCase struct {
	repo       repositories.AssetCommentRepository
	moderation *moderation.Client
}

func NewCreateAssetCommentUseCase(repo repositories.AssetCommentRepository, moderationClient *moderation.Client) *CreateAssetCommentUseCase {
	return &CreateAssetCommentUseCase{repo: repo, moderation: moderationClient}
}

func (uc *CreateAssetCommentUseCase) Execute(ctx context.Context, assetID string, userID string, req request.CreateAssetCommentRequest) (response.AssetCommentResponse, error) {
	if err := req.Validate(); err != nil {
		return response.AssetCommentResponse{}, err
	}

	commentID := uuid.NewString()

	result, err := uc.moderation.Analyze(ctx, commentID, "asset_comment", req.Content)
	if err != nil {
		return response.AssetCommentResponse{}, err
	}
	if result.IsToxic {
		return response.AssetCommentResponse{}, moderation.ErrToxicContent
	}

	created, err := uc.repo.Create(ctx, entities.AssetComment{
		ID:            commentID,
		AssetID:       assetID,
		UserID:        userID,
		Content:       req.Content,
		ToxicityScore: &result.ToxicityScore,
		IsVisible:     true,
	})
	if err != nil {
		return response.AssetCommentResponse{}, err
	}

	return response.FromEntity(created), nil
}
