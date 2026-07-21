package infrastructure

import (
	"vault-payment/src/core/stripeclient"
	"vault-payment/src/features/connect/application"
	"vault-payment/src/features/connect/domain/repositories"
	"vault-payment/src/features/connect/infrastructure/controllers"
)

func BuildCreateOnboardingLinkController(accountRepo repositories.ConnectedAccountRepository, stripeClient stripeclient.Client) *controllers.CreateOnboardingLinkController {
	useCase := application.NewCreateOnboardingLinkUseCase(accountRepo, stripeClient)
	return controllers.NewCreateOnboardingLinkController(useCase)
}

func BuildGetAccountStatusController(accountRepo repositories.ConnectedAccountRepository, stripeClient stripeclient.Client) *controllers.GetAccountStatusController {
	useCase := application.NewGetAccountStatusUseCase(accountRepo, stripeClient)
	return controllers.NewGetAccountStatusController(useCase)
}
