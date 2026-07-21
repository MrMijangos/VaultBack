package response

import "vault-payment/src/features/subscriptions/domain/entities"

type PlanResponse struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	PriceMXN       float64  `json:"price_mxn"`
	MaxAds         int      `json:"max_ads"`
	TargetSections []string `json:"target_sections"`
	CommissionRate float64  `json:"commission_rate"`
}

func PlanFromEntity(p *entities.Plan) PlanResponse {
	return PlanResponse{
		ID:             p.ID,
		Name:           p.Name,
		PriceMXN:       p.PriceMXN,
		MaxAds:         p.MaxAds,
		TargetSections: p.TargetSections,
		CommissionRate: p.CommissionRate,
	}
}
