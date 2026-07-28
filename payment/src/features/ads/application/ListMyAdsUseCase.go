package application

import (
	"context"

	"vault-payment/src/features/ads/domain/dto/response"
	"vault-payment/src/features/ads/domain/repositories"
)

// ListMyAdsUseCase alimenta la pantalla "Mis anuncios" -- a diferencia de
// ListActiveAdsUseCase (público, filtra por sección), esto devuelve TODOS
// los anuncios del dueño (activos e inactivos, ver AdRepository.ListByUserID),
// para que pueda editarlos, borrarlos o ver sus impressions/clicks.
type ListMyAdsUseCase struct {
	adRepo repositories.AdRepository
}

func NewListMyAdsUseCase(adRepo repositories.AdRepository) *ListMyAdsUseCase {
	return &ListMyAdsUseCase{adRepo: adRepo}
}

func (uc *ListMyAdsUseCase) Execute(ctx context.Context, userID string) ([]response.AdResponse, error) {
	ads, err := uc.adRepo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return response.AdsFromEntities(ads), nil
}
