package application

import (
	"context"

	"github.com/google/uuid"

	"vault/src/core/moderation"
	"vault/src/features/assetcomments/domain/dto/request"
	"vault/src/features/assetcomments/domain/dto/response"
	"vault/src/features/assetcomments/domain/entities"
	"vault/src/features/assetcomments/domain/repositories"
	assetrepositories "vault/src/features/assets/domain/repositories"
)

type CreateAssetCommentUseCase struct {
	repo       repositories.AssetCommentRepository
	assetRepo  assetrepositories.AssetRepository
	purchases  repositories.PurchaseVerifier
	moderation *moderation.Client
}

func NewCreateAssetCommentUseCase(
	repo repositories.AssetCommentRepository,
	assetRepo assetrepositories.AssetRepository,
	purchases repositories.PurchaseVerifier,
	moderationClient *moderation.Client,
) *CreateAssetCommentUseCase {
	return &CreateAssetCommentUseCase{repo: repo, assetRepo: assetRepo, purchases: purchases, moderation: moderationClient}
}

func (uc *CreateAssetCommentUseCase) Execute(ctx context.Context, assetID string, userID string, req request.CreateAssetCommentRequest) (response.AssetCommentResponse, error) {
	if err := req.Validate(); err != nil {
		return response.AssetCommentResponse{}, err
	}

	asset, err := uc.assetRepo.FindByID(ctx, assetID)
	if err != nil {
		return response.AssetCommentResponse{}, err
	}

	// Exigir una compra completada antes de gastar la llamada al servicio de
	// moderación -- mismo criterio que reviews.CreateReviewUseCase. El dueño
	// de la publicación queda exento: tiene que poder responder preguntas en
	// su propio artículo sin haberse comprado algo a sí mismo.
	if asset.UserID != userID {
		purchased, err := uc.purchases.HasPurchased(ctx, userID, asset.UserID)
		if err != nil {
			return response.AssetCommentResponse{}, err
		}
		if !purchased {
			return response.AssetCommentResponse{}, ErrNotPurchased
		}
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
