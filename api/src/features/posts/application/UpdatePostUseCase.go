package application

import (
	"context"

	"vault/src/features/posts/domain/dto/request"
	"vault/src/features/posts/domain/dto/response"
	"vault/src/features/posts/domain/repositories"
)

type UpdatePostUseCase struct {
	repo repositories.PostRepository
}

func NewUpdatePostUseCase(repo repositories.PostRepository) *UpdatePostUseCase {
	return &UpdatePostUseCase{repo: repo}
}

func (uc *UpdatePostUseCase) Execute(ctx context.Context, id string, userID string, req request.UpdatePostRequest) (response.PostResponse, error) {
	if err := req.Validate(); err != nil {
		return response.PostResponse{}, err
	}

	updated, err := uc.repo.Update(ctx, id, userID, req.Content)
	if err != nil {
		return response.PostResponse{}, err
	}

	return response.FromEntity(updated, nil), nil
}
