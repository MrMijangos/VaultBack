package repositories

import (
	"context"
	"errors"

	"vault/src/features/posts/domain/entities"
)

var ErrPostNotFound = errors.New("la publicacion no existe")

type PostRepository interface {
	Create(ctx context.Context, post entities.Post) (entities.Post, error)
	FindAllVisible(ctx context.Context) ([]entities.Post, error)
	FindByID(ctx context.Context, id string) (entities.Post, error)
	Update(ctx context.Context, id string, userID string, content string) (entities.Post, error)
	Delete(ctx context.Context, id string, userID string) error
	FindPhotosByPostID(ctx context.Context, postID string) ([]entities.PostPhoto, error)
	AddPhoto(ctx context.Context, postID string, url string) (entities.PostPhoto, error)
	Like(ctx context.Context, postID string, userID string) error
	Unlike(ctx context.Context, postID string, userID string) error
	Save(ctx context.Context, postID string, userID string) error
	Unsave(ctx context.Context, postID string, userID string) error
	FindSavedByUser(ctx context.Context, userID string) ([]entities.Post, error)
}
