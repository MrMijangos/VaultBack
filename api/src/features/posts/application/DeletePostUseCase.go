package application

import (
	"context"

	"vault/src/features/posts/domain/repositories"
)

type DeletePostUseCase struct {
	repo repositories.PostRepository
}

func NewDeletePostUseCase(repo repositories.PostRepository) *DeletePostUseCase {
	return &DeletePostUseCase{repo: repo}
}

func (uc *DeletePostUseCase) Execute(ctx context.Context, id string, userID string) error {
	return uc.repo.Delete(ctx, id, userID)
}
