package application

import (
	"context"

	"vault/src/features/posts/domain/dto/response"
	"vault/src/features/posts/domain/repositories"
)

type GetPostByIdUseCase struct {
	repo repositories.PostRepository
}

func NewGetPostByIdUseCase(repo repositories.PostRepository) *GetPostByIdUseCase {
	return &GetPostByIdUseCase{repo: repo}
}

func (uc *GetPostByIdUseCase) Execute(ctx context.Context, id string) (response.PostResponse, error) {
	p, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return response.PostResponse{}, err
	}

	photos, err := uc.repo.FindPhotosByPostID(ctx, id)
	if err != nil {
		return response.PostResponse{}, err
	}

	return response.FromEntity(p, photos), nil
}
