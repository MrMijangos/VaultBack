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

// Execute trae las fotos de cada post por separado -- igual que en assets,
// las fotos viven en su propia tabla y FromEntities no las incluía, así
// que el feed siempre mostraba los posts sin imagen aunque tuvieran una.
func (uc *GetAllPostsUseCase) Execute(ctx context.Context) ([]response.PostResponse, error) {
	list, err := uc.repo.FindAllVisible(ctx)
	if err != nil {
		return nil, err
	}

	out := make([]response.PostResponse, 0, len(list))
	for _, p := range list {
		photos, err := uc.repo.FindPhotosByPostID(ctx, p.ID)
		if err != nil {
			return nil, err
		}
		out = append(out, response.FromEntity(p, photos))
	}
	return out, nil
}
