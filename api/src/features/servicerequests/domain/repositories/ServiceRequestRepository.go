package repositories

import (
	"context"
	"errors"

	"vault/src/features/servicerequests/domain/entities"
)

var ErrServiceRequestNotFound = errors.New("la solicitud de servicio no existe")

type ServiceRequestRepository interface {
	Create(ctx context.Context, sr *entities.ServiceRequest) error
	Update(ctx context.Context, sr *entities.ServiceRequest) error
	FindByID(ctx context.Context, id string) (*entities.ServiceRequest, error)
	ListByOwnerID(ctx context.Context, ownerID string) ([]entities.ServiceRequest, error)
	ListByBusinessID(ctx context.Context, businessID string) ([]entities.ServiceRequest, error)
}
