package application

import (
	"context"

	"vault-payment/src/features/ads/domain/repositories"
)

type RegisterImpressionUseCase struct {
	adRepo repositories.AdRepository
}

func NewRegisterImpressionUseCase(adRepo repositories.AdRepository) *RegisterImpressionUseCase {
	return &RegisterImpressionUseCase{adRepo: adRepo}
}

func (uc *RegisterImpressionUseCase) Execute(ctx context.Context, adID string) error {
	return uc.adRepo.IncrementImpressions(ctx, adID)
}
