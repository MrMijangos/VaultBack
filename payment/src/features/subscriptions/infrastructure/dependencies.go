package infrastructure

import (
	"vault-payment/src/core/eventbus"
	"vault-payment/src/core/stripeclient"
	"vault-payment/src/features/subscriptions/application"
	"vault-payment/src/features/subscriptions/domain/repositories"
	"vault-payment/src/features/subscriptions/infrastructure/controllers"
)

// Los repos y clientes se construyen una sola vez en main.go (guardan
// estado en memoria) y se inyectan aquí -- a diferencia de api/, donde cada
// Build*Controller crea su propio adapter porque solo envuelve un
// *pgxpool.Pool compartido.

func BuildListPlansController(planRepo repositories.PlanRepository) *controllers.ListPlansController {
	useCase := application.NewListPlansUseCase(planRepo)
	return controllers.NewListPlansController(useCase)
}

func BuildCreateSubscriptionController(
	planRepo repositories.PlanRepository,
	subscriptionRepo repositories.SubscriptionRepository,
	stripeClient stripeclient.Client,
	publisher eventbus.Publisher,
) *controllers.CreateSubscriptionController {
	useCase := application.NewCreateSubscriptionUseCase(planRepo, subscriptionRepo, stripeClient, publisher)
	return controllers.NewCreateSubscriptionController(useCase)
}

func BuildGetSubscriptionStatusController(subscriptionRepo repositories.SubscriptionRepository) *controllers.GetSubscriptionStatusController {
	useCase := application.NewGetSubscriptionStatusUseCase(subscriptionRepo)
	return controllers.NewGetSubscriptionStatusController(useCase)
}

func BuildCancelSubscriptionController(
	subscriptionRepo repositories.SubscriptionRepository,
	stripeClient stripeclient.Client,
	adDeactivator repositories.AdDeactivator,
	publisher eventbus.Publisher,
) *controllers.CancelSubscriptionController {
	useCase := application.NewCancelSubscriptionUseCase(subscriptionRepo, stripeClient, adDeactivator, publisher)
	return controllers.NewCancelSubscriptionController(useCase)
}

func BuildStripeWebhookController(
	subscriptionRepo repositories.SubscriptionRepository,
	stripeClient stripeclient.Client,
	adDeactivator repositories.AdDeactivator,
	publisher eventbus.Publisher,
	webhookSecret string,
) *controllers.StripeWebhookController {
	useCase := application.NewHandleStripeWebhookUseCase(subscriptionRepo, stripeClient, adDeactivator, publisher, webhookSecret)
	return controllers.NewStripeWebhookController(useCase)
}
