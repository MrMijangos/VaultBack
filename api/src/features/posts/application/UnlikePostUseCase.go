package application

import (
	"context"

	"vault/src/features/posts/domain/repositories"
)

type UnlikePostUseCase struct {
	repo repositories.PostRepository
}

func NewUnlikePostUseCase(repo repositories.PostRepository) *UnlikePostUseCase {
	return &UnlikePostUseCase{repo: repo}
}

func (uc *UnlikePostUseCase) Execute(ctx context.Context, postID string, userID string) error {
	return uc.repo.Unlike(ctx, postID, userID)
}
