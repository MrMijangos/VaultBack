package application

import (
	"context"
	"slices"

	"vault-payment/src/features/ads/domain/dto/request"
	"vault-payment/src/features/ads/domain/dto/response"
	"vault-payment/src/features/ads/domain/entities"
	"vault-payment/src/features/ads/domain/repositories"
)

type UpdateAdUseCase struct {
	adRepo repositories.AdRepository
}

func NewUpdateAdUseCase(adRepo repositories.AdRepository) *UpdateAdUseCase {
	return &UpdateAdUseCase{adRepo: adRepo}
}

func (uc *UpdateAdUseCase) Execute(ctx context.Context, userID, adID string, req request.UpdateAdRequest) (*response.AdResponse, error) {
	if !slices.Contains(entities.ValidSections, req.TargetSection) {
		return nil, ErrInvalidSection
	}

	ad, err := uc.adRepo.GetByID(ctx, adID)
	if err != nil {
		return nil, err
	}
	if ad == nil {
		return nil, ErrAdNotFound
	}
	if ad.UserID != userID {
		return nil, ErrNotOwner
	}

	ad.Title = req.Title
	ad.Description = req.Description
	ad.ImageURL = req.ImageURL
	ad.TargetSection = req.TargetSection
	ad.TargetID = req.TargetID

	if err := uc.adRepo.Update(ctx, ad); err != nil {
		return nil, err
	}

	out := response.AdFromEntity(ad)
	return &out, nil
}
