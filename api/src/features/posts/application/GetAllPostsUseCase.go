package application

import (
	"context"

	"vault/src/features/posts/domain/dto/response"
	"vault/src/features/posts/domain/repositories"
)

type GetAllPostsUseCase struct {
	repo repositories.PostRepository
}

func NewGetAllPostsUseCase(repo repositories.PostRepository) *GetAllPostsUseCase {
	return &GetAllPostsUseCase{repo: repo}
}

func (uc *GetAllPostsUseCase) Execute(ctx context.Context) ([]response.PostResponse, error) {
	list, err := uc.repo.FindAllVisible(ctx)
	if err != nil {
		return nil, err
	}
	return response.FromEntities(list), nil
}
