package application

import (
	"context"

	"vault/src/features/comments/domain/repositories"
)

type DeleteCommentUseCase struct {
	repo repositories.CommentRepository
}

func NewDeleteCommentUseCase(repo repositories.CommentRepository) *DeleteCommentUseCase {
	return &DeleteCommentUseCase{repo: repo}
}

func (uc *DeleteCommentUseCase) Execute(ctx context.Context, id string, userID string) error {
	return uc.repo.Delete(ctx, id, userID)
}
