package application

import (
	"context"

	"github.com/google/uuid"

	"vault/src/core/moderation"
	"vault/src/features/posts/domain/dto/request"
	"vault/src/features/posts/domain/dto/response"
	"vault/src/features/posts/domain/entities"
	"vault/src/features/posts/domain/repositories"
)

type CreatePostUseCase struct {
	repo       repositories.PostRepository
	moderation *moderation.Client
}

func NewCreatePostUseCase(repo repositories.PostRepository, moderationClient *moderation.Client) *CreatePostUseCase {
	return &CreatePostUseCase{repo: repo, moderation: moderationClient}
}

func (uc *CreatePostUseCase) Execute(ctx context.Context, userID string, req request.CreatePostRequest) (response.PostResponse, error) {
	if err := req.Validate(); err != nil {
		return response.PostResponse{}, err
	}

	var assetID *string
	if req.AssetID != "" {
		assetID = &req.AssetID
	}

	postID := uuid.NewString()

	result, err := uc.moderation.Analyze(ctx, postID, "post", req.Content)
	if err != nil {
		return response.PostResponse{}, err
	}
	if result.IsToxic {
		return response.PostResponse{}, moderation.ErrToxicContent
	}

	created, err := uc.repo.Create(ctx, entities.Post{
		ID:             postID,
		UserID:         userID,
		AssetID:        assetID,
		Content:        req.Content,
		SentimentScore: &result.SentimentScore,
		SentimentLabel: result.SentimentLabel,
		ToxicityScore:  &result.ToxicityScore,
		IsVisible:      true,
	})
	if err != nil {
		return response.PostResponse{}, err
	}

	return response.FromEntity(created, nil), nil
}
