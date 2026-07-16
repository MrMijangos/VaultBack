package application

import (
	"context"

	"github.com/google/uuid"

	"vault/src/core/moderation"
	"vault/src/features/reviews/domain/dto/request"
	"vault/src/features/reviews/domain/dto/response"
	"vault/src/features/reviews/domain/entities"
	"vault/src/features/reviews/domain/repositories"
)

type CreateReviewUseCase struct {
	repo       repositories.ReviewRepository
	moderation *moderation.Client
}

func NewCreateReviewUseCase(repo repositories.ReviewRepository, moderationClient *moderation.Client) *CreateReviewUseCase {
	return &CreateReviewUseCase{repo: repo, moderation: moderationClient}
}

func (uc *CreateReviewUseCase) Execute(ctx context.Context, userID string, req request.CreateReviewRequest) (response.ReviewResponse, error) {
	if err := req.Validate(); err != nil {
		return response.ReviewResponse{}, err
	}

	reviewID := uuid.NewString()

	result, err := uc.moderation.Analyze(ctx, reviewID, "review", req.Content)
	if err != nil {
		return response.ReviewResponse{}, err
	}
	if result.IsToxic {
		return response.ReviewResponse{}, moderation.ErrToxicContent
	}

	created, err := uc.repo.Create(ctx, entities.Review{
		ID:             reviewID,
		UserID:         userID,
		ProviderID:     req.ProviderID,
		Content:        req.Content,
		SentimentScore: &result.SentimentScore,
		ToxicityScore:  &result.ToxicityScore,
		IsVisible:      true,
	})
	if err != nil {
		return response.ReviewResponse{}, err
	}

	return response.FromEntity(created), nil
}
