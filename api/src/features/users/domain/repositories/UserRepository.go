package repositories

import (
	"context"
	"errors"

	"vault/src/features/users/domain/entities"
)

var ErrUserNotFound = errors.New("el usuario no existe")

type UserRepository interface {
	Create(ctx context.Context, user entities.User) (entities.User, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	FindAll(ctx context.Context) ([]entities.User, error)
	FindByID(ctx context.Context, id string) (entities.User, error)
	Update(ctx context.Context, id string, user entities.User) (entities.User, error)
	UpdateImage(ctx context.Context, id string, imageURL string) (entities.User, error)
	Delete(ctx context.Context, id string) error
}
