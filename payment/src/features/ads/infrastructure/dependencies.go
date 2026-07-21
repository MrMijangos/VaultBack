package infrastructure

import (
	"vault-payment/src/features/ads/application"
	"vault-payment/src/features/ads/domain/repositories"
	"vault-payment/src/features/ads/infrastructure/controllers"
)

func BuildCreateAdController(adRepo repositories.AdRepository, subscriptionProvider repositories.SubscriptionInfoProvider) *controllers.CreateAdController {
	useCase := application.NewCreateAdUseCase(adRepo, subscriptionProvider)
	return controllers.NewCreateAdController(useCase)
}

func BuildUpdateAdController(adRepo repositories.AdRepository) *controllers.UpdateAdController {
	useCase := application.NewUpdateAdUseCase(adRepo)
	return controllers.NewUpdateAdController(useCase)
}

func BuildDeleteAdController(adRepo repositories.AdRepository) *controllers.DeleteAdController {
	useCase := application.NewDeleteAdUseCase(adRepo)
	return controllers.NewDeleteAdController(useCase)
}

func BuildListActiveAdsController(adRepo repositories.AdRepository) *controllers.ListActiveAdsController {
	useCase := application.NewListActiveAdsUseCase(adRepo)
	return controllers.NewListActiveAdsController(useCase)
}

func BuildRegisterImpressionController(adRepo repositories.AdRepository) *controllers.RegisterImpressionController {
	useCase := application.NewRegisterImpressionUseCase(adRepo)
	return controllers.NewRegisterImpressionController(useCase)
}

func BuildRegisterClickController(adRepo repositories.AdRepository) *controllers.RegisterClickController {
	useCase := application.NewRegisterClickUseCase(adRepo)
	return controllers.NewRegisterClickController(useCase)
}
