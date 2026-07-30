package application

import (
	"context"
	"errors"

	businessrepositories "vault/src/features/businesses/domain/repositories"
	"vault/src/features/servicerequests/domain/dto/response"
	"vault/src/features/servicerequests/domain/repositories"
)

// ListIncomingServiceRequestsUseCase lista las solicitudes que le llegaron
// al negocio del usuario autenticado (para "Mis Negocios"). Si el usuario
// no tiene negocio, devuelve una lista vacía en vez de error -- no tiene
// nada que ver todavía, no es una condición de falla.
type ListIncomingServiceRequestsUseCase struct {
	repo         repositories.ServiceRequestRepository
	businessRepo businessrepositories.BusinessRepository
}

func NewListIncomingServiceRequestsUseCase(
	repo repositories.ServiceRequestRepository,
	businessRepo businessrepositories.BusinessRepository,
) *ListIncomingServiceRequestsUseCase {
	return &ListIncomingServiceRequestsUseCase{repo: repo, businessRepo: businessRepo}
}

func (uc *ListIncomingServiceRequestsUseCase) Execute(ctx context.Context, userID string) ([]response.ServiceRequestResponse, error) {
	business, err := uc.businessRepo.FindByUserID(ctx, userID)
	if errors.Is(err, businessrepositories.ErrBusinessNotFound) {
		return []response.ServiceRequestResponse{}, nil
	}
	if err != nil {
		return nil, err
	}

	list, err := uc.repo.ListByBusinessID(ctx, business.ID)
	if err != nil {
		return nil, err
	}
	return response.FromEntities(list), nil
}
