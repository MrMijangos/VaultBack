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
//
// userID viene vacío para un visitante sin sesión (la ruta es pública,
// ver OptionalAuth) -- en ese caso is_liked/is_saved quedan en false para
// todos, no se consulta nada de más.
func (uc *GetAllPostsUseCase) Execute(ctx context.Context, userID string) ([]response.PostResponse, error) {
	list, err := uc.repo.FindAllVisible(ctx)
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
	saved, err := uc.repo.FindSavedPostIDs(ctx, postIDs, userID)
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
		item.IsSaved = saved[p.ID]
		out = append(out, item)
	}
	return out, nil
}
