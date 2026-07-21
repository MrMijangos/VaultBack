package adapters

import (
	"context"
	"fmt"

	"vault-payment/src/core/config"
	"vault-payment/src/features/subscriptions/domain/entities"
)

// InMemoryPlanRepository sirve los 3 planes fijos definidos en la
// especificación. Es de solo lectura -- los planes no se crean/editan por
// API, solo se seedean al arrancar. El StripePriceID se toma de las
// variables de entorno para no hardcodear IDs de Stripe en el código (útil
// también porque todavía no hay cuenta de Stripe: pueden quedar vacíos).
type InMemoryPlanRepository struct {
	plans []*entities.Plan
}

func NewInMemoryPlanRepository(cfg *config.Config) *InMemoryPlanRepository {
	return &InMemoryPlanRepository{
		plans: []*entities.Plan{
			{
				ID:             entities.PlanBasico,
				Name:           "Básico",
				PriceMXN:       49,
				StripePriceID:  cfg.StripePriceBasico,
				MaxAds:         1,
				TargetSections: []string{entities.SectionMarketplace},
				CommissionRate: 0.08,
			},
			{
				ID:             entities.PlanPro,
				Name:           "Pro",
				PriceMXN:       99,
				StripePriceID:  cfg.StripePricePro,
				MaxAds:         3,
				TargetSections: []string{entities.SectionMarketplace, entities.SectionFeed},
				CommissionRate: 0.05,
			},
			{
				ID:             entities.PlanPremium,
				Name:           "Premium",
				PriceMXN:       179,
				StripePriceID:  cfg.StripePricePremium,
				MaxAds:         5,
				TargetSections: []string{entities.SectionMarketplace, entities.SectionFeed},
				CommissionRate: 0.03,
			},
		},
	}
}

func (r *InMemoryPlanRepository) List(ctx context.Context) ([]*entities.Plan, error) {
	return r.plans, nil
}

func (r *InMemoryPlanRepository) GetByID(ctx context.Context, id string) (*entities.Plan, error) {
	for _, p := range r.plans {
		if p.ID == id {
			return p, nil
		}
	}
	return nil, fmt.Errorf("plan %q no existe", id)
}
