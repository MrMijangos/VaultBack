package repositories

import (
	"context"
	"errors"

	"vault/src/features/businessservices/domain/entities"
)

var ErrBusinessServiceNotFound = errors.New("el servicio no existe")

type BusinessServiceRepository interface {
	Create(ctx context.Context, service entities.BusinessService) (entities.BusinessService, error)
	Update(ctx context.Context, id string, businessID string, service entities.BusinessService) (entities.BusinessService, error)
	Delete(ctx context.Context, id string, businessID string) error
	ListByBusinessID(ctx context.Context, businessID string) ([]entities.BusinessService, error)
}
