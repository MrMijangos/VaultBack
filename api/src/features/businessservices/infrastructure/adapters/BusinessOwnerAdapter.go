package adapters

import (
	"context"
	"errors"

	businessesRepositories "vault/src/features/businesses/domain/repositories"
	businessesAdapters "vault/src/features/businesses/infrastructure/adapters"

	"vault/src/features/businessservices/domain/repositories"
)

// BusinessOwnerAdapter satisface repositories.BusinessOwnerProvider apoyándose
// en el repositorio real de businesses, sin que el dominio de
// businessservices dependa de sus tipos concretos.
type BusinessOwnerAdapter struct {
	repo *businessesAdapters.PostgreSQLBusinessRepository
}

func NewBusinessOwnerAdapter(repo *businessesAdapters.PostgreSQLBusinessRepository) *BusinessOwnerAdapter {
	return &BusinessOwnerAdapter{repo: repo}
}

func (a *BusinessOwnerAdapter) GetOwnerUserID(ctx context.Context, businessID string) (string, error) {
	business, err := a.repo.FindByID(ctx, businessID)
	if errors.Is(err, businessesRepositories.ErrBusinessNotFound) {
		return "", repositories.ErrBusinessNotFound
	}
	if err != nil {
		return "", err
	}
	return business.UserID, nil
}
