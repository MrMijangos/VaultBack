package application

import (
	"context"

	"vault/src/features/posts/domain/dto/request"
	"vault/src/features/posts/domain/dto/response"
	"vault/src/features/posts/domain/entities"
	"vault/src/features/posts/domain/repositories"
)

type CreatePostUseCase struct {
	repo repositories.PostRepository
}

func NewCreatePostUseCase(repo repositories.PostRepository) *CreatePostUseCase {
	return &CreatePostUseCase{repo: repo}
}

func (uc *CreatePostUseCase) Execute(ctx context.Context, userID string, req request.CreatePostRequest) (response.PostResponse, error) {
	if err := req.Validate(); err != nil {
		return response.PostResponse{}, err
	}

	var assetID *string
	if req.AssetID != "" {
		assetID = &req.AssetID
	}

	created, err := uc.repo.Create(ctx, entities.Post{
		UserID:  userID,
		AssetID: assetID,
		Content: req.Content,
	})
	if err != nil {
		return response.PostResponse{}, err
	}

	return response.FromEntity(created, nil), nil
}
