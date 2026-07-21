package application

import (
	"context"

	"vault-payment/src/features/ads/domain/repositories"
)

type RegisterClickUseCase struct {
	adRepo repositories.AdRepository
}

func NewRegisterClickUseCase(adRepo repositories.AdRepository) *RegisterClickUseCase {
	return &RegisterClickUseCase{adRepo: adRepo}
}

func (uc *RegisterClickUseCase) Execute(ctx context.Context, adID string) error {
	return uc.adRepo.IncrementClicks(ctx, adID)
}
