package application

import (
	"context"

	"vault/src/features/businessservices/domain/repositories"
)

type DeleteBusinessServiceUseCase struct {
	repo          repositories.BusinessServiceRepository
	ownerProvider repositories.BusinessOwnerProvider
}

func NewDeleteBusinessServiceUseCase(repo repositories.BusinessServiceRepository, ownerProvider repositories.BusinessOwnerProvider) *DeleteBusinessServiceUseCase {
	return &DeleteBusinessServiceUseCase{repo: repo, ownerProvider: ownerProvider}
}

func (uc *DeleteBusinessServiceUseCase) Execute(ctx context.Context, businessID string, serviceID string, userID string) error {
	ownerID, err := uc.ownerProvider.GetOwnerUserID(ctx, businessID)
	if err != nil {
		return err
	}
	if ownerID != userID {
		return repositories.ErrNotOwner
	}

	return uc.repo.Delete(ctx, serviceID, businessID)
}
