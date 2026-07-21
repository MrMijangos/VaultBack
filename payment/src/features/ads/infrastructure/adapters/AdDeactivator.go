package adapters

import (
	"context"
	"fmt"

	"vault-payment/src/features/ads/domain/entities"
	"vault-payment/src/features/ads/domain/repositories"
)

// AdDeactivator implementa subscriptions/domain/repositories.AdDeactivator
// (tipado estructural, sin import cruzado entre features).
type AdDeactivator struct {
	repo repositories.AdRepository
}

func NewAdDeactivator(repo repositories.AdRepository) *AdDeactivator {
	return &AdDeactivator{repo: repo}
}

func (d *AdDeactivator) DeactivateBySubscriptionID(ctx context.Context, subscriptionID string) error {
	ads, err := d.repo.ListBySubscriptionID(ctx, subscriptionID)
	if err != nil {
		return fmt.Errorf("no se pudieron listar los anuncios de la suscripción %q: %w", subscriptionID, err)
	}

	for _, ad := range ads {
		if ad.Status == entities.AdStatusInactive {
			continue
		}
		ad.Status = entities.AdStatusInactive
		if err := d.repo.Update(ctx, ad); err != nil {
			return fmt.Errorf("no se pudo desactivar el anuncio %q: %w", ad.ID, err)
		}
	}
	return nil
}
