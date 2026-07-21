package application

import (
	"context"
	"slices"
	"time"

	"github.com/google/uuid"

	"vault-payment/src/features/ads/domain/dto/request"
	"vault-payment/src/features/ads/domain/dto/response"
	"vault-payment/src/features/ads/domain/entities"
	"vault-payment/src/features/ads/domain/repositories"
)

type CreateAdUseCase struct {
	adRepo               repositories.AdRepository
	subscriptionProvider repositories.SubscriptionInfoProvider
}

func NewCreateAdUseCase(adRepo repositories.AdRepository, subscriptionProvider repositories.SubscriptionInfoProvider) *CreateAdUseCase {
	return &CreateAdUseCase{adRepo: adRepo, subscriptionProvider: subscriptionProvider}
}

func (uc *CreateAdUseCase) Execute(ctx context.Context, userID string, req request.CreateAdRequest) (*response.AdResponse, error) {
	if !slices.Contains(entities.ValidSections, req.TargetSection) {
		return nil, ErrInvalidSection
	}

	subscriptionID, maxAds, targetSections, found, err := uc.subscriptionProvider.GetActiveSubscription(ctx, userID)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, ErrNoActiveSubscription
	}
	if !slices.Contains(targetSections, req.TargetSection) {
		return nil, ErrSectionNotAllowed
	}

	activeCount, err := uc.adRepo.CountActiveByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if activeCount >= maxAds {
		return nil, ErrMaxAdsReached
	}

	ad := &entities.Ad{
		ID:             uuid.NewString(),
		UserID:         userID,
		SubscriptionID: subscriptionID,
		Title:          req.Title,
		Description:    req.Description,
		ImageURL:       req.ImageURL,
		TargetSection:  req.TargetSection,
		TargetID:       req.TargetID,
		Status:         entities.AdStatusActive,
		CreatedAt:      time.Now().UTC(),
	}

	if err := uc.adRepo.Create(ctx, ad); err != nil {
		return nil, err
	}

	out := response.AdFromEntity(ad)
	return &out, nil
}
