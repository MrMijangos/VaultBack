package application

import (
	"context"

	"vault/src/features/posts/domain/repositories"
)

type UnsavePostUseCase struct {
	repo repositories.PostRepository
}

func NewUnsavePostUseCase(repo repositories.PostRepository) *UnsavePostUseCase {
	return &UnsavePostUseCase{repo: repo}
}

func (uc *UnsavePostUseCase) Execute(ctx context.Context, postID string, userID string) error {
	return uc.repo.Unsave(ctx, postID, userID)
}
