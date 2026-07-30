package application

import (
	"context"

	businessrepositories "vault/src/features/businesses/domain/repositories"
	"vault/src/features/servicerequests/domain/entities"
	"vault/src/features/servicerequests/domain/repositories"
)

// loadForBusinessAction busca la solicitud, confirma que businessOwnerID es
// dueño del negocio al que pertenece, y que su estado actual es
// expectedStatus -- lo comparten Accept/Start/Finish, que son exactamente
// la misma forma (cargar, autorizar, validar transición) con distinto
// destino.
func loadForBusinessAction(
	ctx context.Context,
	repo repositories.ServiceRequestRepository,
	businessRepo businessrepositories.BusinessRepository,
	businessOwnerID string,
	requestID string,
	expectedStatus string,
) (*entities.ServiceRequest, error) {
	sr, err := repo.FindByID(ctx, requestID)
	if err != nil {
		return nil, err
	}
	if sr == nil {
		return nil, repositories.ErrServiceRequestNotFound
	}

	business, err := businessRepo.FindByID(ctx, sr.BusinessID)
	if err != nil {
		return nil, err
	}
	if business.UserID != businessOwnerID {
		return nil, ErrNotBusinessOwner
	}
	if sr.Status != expectedStatus {
		return nil, ErrInvalidTransition
	}
	return sr, nil
}
