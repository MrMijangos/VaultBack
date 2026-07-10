package infrastructure

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"vault/src/features/notifications/application"
	"vault/src/features/notifications/infrastructure/adapters"
	"vault/src/features/notifications/infrastructure/controllers"
)

func BuildCreateNotificationController(pool *pgxpool.Pool) *controllers.CreateNotificationController {
	repo := adapters.NewPostgreSQLNotificationRepository(pool)
	useCase := application.NewCreateNotificationUseCase(repo)
	return controllers.NewCreateNotificationController(useCase)
}

func BuildGetMyNotificationsController(pool *pgxpool.Pool) *controllers.GetMyNotificationsController {
	repo := adapters.NewPostgreSQLNotificationRepository(pool)
	useCase := application.NewGetMyNotificationsUseCase(repo)
	return controllers.NewGetMyNotificationsController(useCase)
}

func BuildMarkNotificationAsReadController(pool *pgxpool.Pool) *controllers.MarkNotificationAsReadController {
	repo := adapters.NewPostgreSQLNotificationRepository(pool)
	useCase := application.NewMarkNotificationAsReadUseCase(repo)
	return controllers.NewMarkNotificationAsReadController(useCase)
}

func BuildDeleteNotificationController(pool *pgxpool.Pool) *controllers.DeleteNotificationController {
	repo := adapters.NewPostgreSQLNotificationRepository(pool)
	useCase := application.NewDeleteNotificationUseCase(repo)
	return controllers.NewDeleteNotificationController(useCase)
}
