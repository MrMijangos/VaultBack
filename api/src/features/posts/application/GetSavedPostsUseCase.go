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

	postIDs := make([]string, len(list))
	for i, p := range list {
		postIDs[i] = p.ID
	}
	liked, err := uc.repo.FindLikedPostIDs(ctx, postIDs, userID)
	if err != nil {
		return nil, err
	}

	out := make([]response.PostResponse, 0, len(list))
	for _, p := range list {
		photos, err := uc.repo.FindPhotosByPostID(ctx, p.ID)
		if err != nil {
			return nil, err
		}
		item := response.FromEntity(p, photos)
		item.IsLiked = liked[p.ID]
		item.IsSaved = true // por definición: esta lista solo trae guardados.
		out = append(out, item)
	}
	return out, nil
}
