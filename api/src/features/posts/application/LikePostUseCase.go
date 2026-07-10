package application

import (
	"context"

	"vault/src/features/posts/domain/repositories"
)

type LikePostUseCase struct {
	repo repositories.PostRepository
}

func NewLikePostUseCase(repo repositories.PostRepository) *LikePostUseCase {
	return &LikePostUseCase{repo: repo}
}

func (uc *LikePostUseCase) Execute(ctx context.Context, postID string, userID string) error {
	return uc.repo.Like(ctx, postID, userID)
}
