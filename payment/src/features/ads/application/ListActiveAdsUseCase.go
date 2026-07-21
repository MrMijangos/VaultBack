package application

import (
	"context"

	"vault-payment/src/features/ads/domain/dto/response"
	"vault-payment/src/features/ads/domain/repositories"
)

// ListActiveAdsUseCase alimenta la mezcla de anuncios en el feed/marketplace
// del lado de Flutter (todavía no implementado -- se hace cuando el usuario
// lo pida). El endpoint ya queda listo del lado del backend.
type ListActiveAdsUseCase struct {
	adRepo repositories.AdRepository
}

func NewListActiveAdsUseCase(adRepo repositories.AdRepository) *ListActiveAdsUseCase {
	return &ListActiveAdsUseCase{adRepo: adRepo}
}

func (uc *ListActiveAdsUseCase) Execute(ctx context.Context, section string) ([]response.AdResponse, error) {
	ads, err := uc.adRepo.ListActiveBySection(ctx, section)
	if err != nil {
		return nil, err
	}
	return response.AdsFromEntities(ads), nil
}
