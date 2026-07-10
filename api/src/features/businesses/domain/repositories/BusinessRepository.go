package repositories

import (
	"context"
	"errors"

	"vault/src/features/businesses/domain/entities"
)

var ErrBusinessNotFound = errors.New("el negocio no existe")
var ErrBusinessAlreadyExists = errors.New("ya tienes un negocio registrado")

type BusinessRepository interface {
	Create(ctx context.Context, business entities.Business) (entities.Business, error)
	ExistsByUserID(ctx context.Context, userID string) (bool, error)
	FindAll(ctx context.Context) ([]entities.Business, error)
	FindByID(ctx context.Context, id string) (entities.Business, error)
	Update(ctx context.Context, id string, userID string, business entities.Business) (entities.Business, error)
	Delete(ctx context.Context, id string, userID string) error
}
