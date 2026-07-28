package application

import (
	"context"

	"github.com/google/uuid"

	"vault/src/core/moderation"
	"vault/src/features/posts/domain/dto/request"
	"vault/src/features/posts/domain/dto/response"
	"vault/src/features/posts/domain/entities"
	"vault/src/features/posts/domain/repositories"
)

type CreatePostUseCase struct {
	repo        repositories.PostRepository
	moderation  *moderation.Client
	assetPhotos repositories.AssetPhotoProvider
}

func NewCreatePostUseCase(repo repositories.PostRepository, moderationClient *moderation.Client, assetPhotos repositories.AssetPhotoProvider) *CreatePostUseCase {
	return &CreatePostUseCase{repo: repo, moderation: moderationClient, assetPhotos: assetPhotos}
}

func (uc *CreatePostUseCase) Execute(ctx context.Context, userID string, req request.CreatePostRequest) (response.PostResponse, error) {
	if err := req.Validate(); err != nil {
		return response.PostResponse{}, err
	}

	var assetID *string
	if req.AssetID != "" {
		assetID = &req.AssetID
	}

	postID := uuid.NewString()

	result, err := uc.moderation.Analyze(ctx, postID, "post", req.Content)
	if err != nil {
		return response.PostResponse{}, err
	}
	if result.IsToxic {
		return response.PostResponse{}, moderation.ErrToxicContent
	}

	created, err := uc.repo.Create(ctx, entities.Post{
		ID:             postID,
		UserID:         userID,
		AssetID:        assetID,
		Content:        req.Content,
		SentimentScore: &result.SentimentScore,
		SentimentLabel: result.SentimentLabel,
		ToxicityScore:  &result.ToxicityScore,
		IsVisible:      true,
	})
	if err != nil {
		return response.PostResponse{}, err
	}

	// "Publicar en el Feed" desde un activo manda asset_id sin ninguna foto
	// propia (vault-app solo re-sube fotos para posts creados desde cero) --
	// sin esto, el post quedaba siempre sin imagen aunque el activo sí
	// tuviera fotos. Se copian las URLs ya existentes en Cloudinary, no hace
	// falta volver a subir el archivo.
	var photos []entities.PostPhoto
	if assetID != nil {
		assetPhotos, err := uc.assetPhotos.FindPhotosByAssetID(ctx, *assetID)
		if err != nil {
			return response.PostResponse{}, err
		}
		for _, ap := range assetPhotos {
			added, err := uc.repo.AddPhoto(ctx, postID, ap.URL)
			if err != nil {
				return response.PostResponse{}, err
			}
			photos = append(photos, added)
		}
	}

	return response.FromEntity(created, photos), nil
}
