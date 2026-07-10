package application

import (
	"context"

	"vault/src/features/comments/domain/dto/response"
	"vault/src/features/comments/domain/repositories"
)

type GetCommentsByPostUseCase struct {
	repo repositories.CommentRepository
}

func NewGetCommentsByPostUseCase(repo repositories.CommentRepository) *GetCommentsByPostUseCase {
	return &GetCommentsByPostUseCase{repo: repo}
}

func (uc *GetCommentsByPostUseCase) Execute(ctx context.Context, postID string) ([]response.CommentResponse, error) {
	list, err := uc.repo.FindByPostID(ctx, postID)
	if err != nil {
		return nil, err
	}
	return response.FromEntities(list), nil
}
