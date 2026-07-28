package infrastructure

import (
	"vault-payment/src/core/eventbus"
	"vault-payment/src/core/stripeclient"
	"vault-payment/src/features/orders/application"
	"vault-payment/src/features/orders/domain/repositories"
	"vault-payment/src/features/orders/infrastructure/controllers"
)

func BuildCreateOrderController(
	orderRepo repositories.OrderRepository,
	commissionProvider repositories.SellerCommissionProvider,
	accountProvider repositories.SellerAccountProvider,
	stripeClient stripeclient.Client,
	publisher eventbus.Publisher,
) *controllers.CreateOrderController {
	useCase := application.NewCreateOrderUseCase(orderRepo, commissionProvider, accountProvider, stripeClient, publisher)
	return controllers.NewCreateOrderController(useCase)
}

func BuildConfirmOrderController(
	orderRepo repositories.OrderRepository,
	accountProvider repositories.SellerAccountProvider,
	stripeClient stripeclient.Client,
	publisher eventbus.Publisher,
) *controllers.ConfirmOrderController {
	useCase := application.NewConfirmOrderUseCase(orderRepo, accountProvider, stripeClient, publisher)
	return controllers.NewConfirmOrderController(useCase)
}

func BuildGetOrderController(orderRepo repositories.OrderRepository) *controllers.GetOrderController {
	useCase := application.NewGetOrderUseCase(orderRepo)
	return controllers.NewGetOrderController(useCase)
}

func BuildShipOrderController(
	orderRepo repositories.OrderRepository,
	publisher eventbus.Publisher,
) *controllers.ShipOrderController {
	useCase := application.NewShipOrderUseCase(orderRepo, publisher)
	return controllers.NewShipOrderController(useCase)
}

func BuildListMyOrdersController(orderRepo repositories.OrderRepository) *controllers.ListMyOrdersController {
	useCase := application.NewListMyOrdersUseCase(orderRepo)
	return controllers.NewListMyOrdersController(useCase)
}

func BuildListMySalesController(orderRepo repositories.OrderRepository) *controllers.ListMySalesController {
	useCase := application.NewListMySalesUseCase(orderRepo)
	return controllers.NewListMySalesController(useCase)
}
