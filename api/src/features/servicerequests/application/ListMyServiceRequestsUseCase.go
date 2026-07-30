package application

import (
	"context"

	"vault/src/features/servicerequests/domain/dto/response"
	"vault/src/features/servicerequests/domain/repositories"
)

// ListMyServiceRequestsUseCase lista las solicitudes que el usuario mandó
// como dueño de un activo (para "Mis Activos": saber cuáles están en curso
// o listas para confirmar).
type ListMyServiceRequestsUseCase struct {
	repo repositories.ServiceRequestRepository
}

func NewListMyServiceRequestsUseCase(repo repositories.ServiceRequestRepository) *ListMyServiceRequestsUseCase {
	return &ListMyServiceRequestsUseCase{repo: repo}
}

func (uc *ListMyServiceRequestsUseCase) Execute(ctx context.Context, ownerID string) ([]response.ServiceRequestResponse, error) {
	list, err := uc.repo.ListByOwnerID(ctx, ownerID)
	if err != nil {
		return nil, err
	}
	return response.FromEntities(list), nil
}
