package repositories

import (
	"context"
	"errors"

	"vault/src/features/comments/domain/entities"
)

var ErrCommentNotFound = errors.New("el comentario no existe")

type CommentRepository interface {
	Create(ctx context.Context, comment entities.Comment) (entities.Comment, error)
	FindByPostID(ctx context.Context, postID string) ([]entities.Comment, error)
	Delete(ctx context.Context, id string, userID string) error
}
