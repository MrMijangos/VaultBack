package application

import (
	"context"

	"vault-payment/src/features/ads/domain/repositories"
)

type DeleteAdUseCase struct {
	adRepo repositories.AdRepository
}

func NewDeleteAdUseCase(adRepo repositories.AdRepository) *DeleteAdUseCase {
	return &DeleteAdUseCase{adRepo: adRepo}
}

func (uc *DeleteAdUseCase) Execute(ctx context.Context, userID, adID string) error {
	ad, err := uc.adRepo.GetByID(ctx, adID)
	if err != nil {
		return err
	}
	if ad == nil {
		return ErrAdNotFound
	}
	if ad.UserID != userID {
		return ErrNotOwner
	}

	return uc.adRepo.Delete(ctx, adID)
}
