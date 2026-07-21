package application

import (
	"context"

	"vault/src/features/posts/domain/repositories"
)

type SavePostUseCase struct {
	repo repositories.PostRepository
}

func NewSavePostUseCase(repo repositories.PostRepository) *SavePostUseCase {
	return &SavePostUseCase{repo: repo}
}

func (uc *SavePostUseCase) Execute(ctx context.Context, postID string, userID string) error {
	return uc.repo.Save(ctx, postID, userID)
}
