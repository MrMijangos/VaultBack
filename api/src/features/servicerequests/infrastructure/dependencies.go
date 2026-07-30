package infrastructure

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"vault/src/core/push"
	assetsAdapters "vault/src/features/assets/infrastructure/adapters"
	businessesAdapters "vault/src/features/businesses/infrastructure/adapters"
	maintenancelogsAdapters "vault/src/features/maintenancelogs/infrastructure/adapters"
	notificationsInfra "vault/src/features/notifications/infrastructure"
	"vault/src/features/servicerequests/application"
	"vault/src/features/servicerequests/infrastructure/adapters"
	"vault/src/features/servicerequests/infrastructure/controllers"
)

func BuildCreateServiceRequestController(pool *pgxpool.Pool, sender push.Sender) *controllers.CreateServiceRequestController {
	repo := adapters.NewPostgreSQLServiceRequestRepository(pool)
	assetRepo := assetsAdapters.NewPostgreSQLAssetRepository(pool)
	businessRepo := businessesAdapters.NewPostgreSQLBusinessRepository(pool)
	notifier := notificationsInfra.BuildNotificationPublisher(pool, sender)
	useCase := application.NewCreateServiceRequestUseCase(repo, assetRepo, businessRepo, notifier)
	return controllers.NewCreateServiceRequestController(useCase)
}

func BuildAcceptServiceRequestController(pool *pgxpool.Pool, sender push.Sender) *controllers.AcceptServiceRequestController {
	repo := adapters.NewPostgreSQLServiceRequestRepository(pool)
	businessRepo := businessesAdapters.NewPostgreSQLBusinessRepository(pool)
	notifier := notificationsInfra.BuildNotificationPublisher(pool, sender)
	useCase := application.NewAcceptServiceRequestUseCase(repo, businessRepo, notifier)
	return controllers.NewAcceptServiceRequestController(useCase)
}

func BuildStartServiceRequestController(pool *pgxpool.Pool, sender push.Sender) *controllers.StartServiceRequestController {
	repo := adapters.NewPostgreSQLServiceRequestRepository(pool)
	businessRepo := businessesAdapters.NewPostgreSQLBusinessRepository(pool)
	notifier := notificationsInfra.BuildNotificationPublisher(pool, sender)
	useCase := application.NewStartServiceRequestUseCase(repo, businessRepo, notifier)
	return controllers.NewStartServiceRequestController(useCase)
}

func BuildFinishServiceRequestController(pool *pgxpool.Pool, sender push.Sender) *controllers.FinishServiceRequestController {
	repo := adapters.NewPostgreSQLServiceRequestRepository(pool)
	businessRepo := businessesAdapters.NewPostgreSQLBusinessRepository(pool)
	notifier := notificationsInfra.BuildNotificationPublisher(pool, sender)
	useCase := application.NewFinishServiceRequestUseCase(repo, businessRepo, notifier)
	return controllers.NewFinishServiceRequestController(useCase)
}

func BuildConfirmServiceRequestController(pool *pgxpool.Pool, sender push.Sender) *controllers.ConfirmServiceRequestController {
	repo := adapters.NewPostgreSQLServiceRequestRepository(pool)
	businessRepo := businessesAdapters.NewPostgreSQLBusinessRepository(pool)
	maintenanceLogRepo := maintenancelogsAdapters.NewPostgreSQLMaintenanceLogRepository(pool)
	notifier := notificationsInfra.BuildNotificationPublisher(pool, sender)
	useCase := application.NewConfirmServiceRequestUseCase(repo, businessRepo, maintenanceLogRepo, notifier)
	return controllers.NewConfirmServiceRequestController(useCase)
}

func BuildListMyServiceRequestsController(pool *pgxpool.Pool) *controllers.ListMyServiceRequestsController {
	repo := adapters.NewPostgreSQLServiceRequestRepository(pool)
	useCase := application.NewListMyServiceRequestsUseCase(repo)
	return controllers.NewListMyServiceRequestsController(useCase)
}

func BuildListIncomingServiceRequestsController(pool *pgxpool.Pool) *controllers.ListIncomingServiceRequestsController {
	repo := adapters.NewPostgreSQLServiceRequestRepository(pool)
	businessRepo := businessesAdapters.NewPostgreSQLBusinessRepository(pool)
	useCase := application.NewListIncomingServiceRequestsUseCase(repo, businessRepo)
	return controllers.NewListIncomingServiceRequestsController(useCase)
}
