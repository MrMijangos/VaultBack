package application

import (
	"context"

	"vault/src/features/posts/domain/dto/response"
	"vault/src/features/posts/domain/repositories"
)

type GetSavedPostsUseCase struct {
	repo repositories.PostRepository
}

func NewGetSavedPostsUseCase(repo repositories.PostRepository) *GetSavedPostsUseCase {
	return &GetSavedPostsUseCase{repo: repo}
}

func (uc *GetSavedPostsUseCase) Execute(ctx context.Context, userID string) ([]response.PostResponse, error) {
	list, err := uc.repo.FindSavedByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	return response.FromEntities(list), nil
}
