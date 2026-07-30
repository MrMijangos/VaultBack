package infrastructure

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"vault/src/core/push"
	"vault/src/features/notifications/application"
	"vault/src/features/notifications/infrastructure/adapters"
	"vault/src/features/notifications/infrastructure/controllers"
)

func BuildCreateNotificationController(pool *pgxpool.Pool, sender push.Sender) *controllers.CreateNotificationController {
	repo := adapters.NewPostgreSQLNotificationRepository(pool, sender)
	useCase := application.NewCreateNotificationUseCase(repo)
	return controllers.NewCreateNotificationController(useCase)
}

func BuildGetMyNotificationsController(pool *pgxpool.Pool, sender push.Sender) *controllers.GetMyNotificationsController {
	repo := adapters.NewPostgreSQLNotificationRepository(pool, sender)
	useCase := application.NewGetMyNotificationsUseCase(repo)
	return controllers.NewGetMyNotificationsController(useCase)
}

func BuildMarkNotificationAsReadController(pool *pgxpool.Pool, sender push.Sender) *controllers.MarkNotificationAsReadController {
	repo := adapters.NewPostgreSQLNotificationRepository(pool, sender)
	useCase := application.NewMarkNotificationAsReadUseCase(repo)
	return controllers.NewMarkNotificationAsReadController(useCase)
}

func BuildMarkAllNotificationsAsReadController(pool *pgxpool.Pool, sender push.Sender) *controllers.MarkAllNotificationsAsReadController {
	repo := adapters.NewPostgreSQLNotificationRepository(pool, sender)
	useCase := application.NewMarkAllNotificationsAsReadUseCase(repo)
	return controllers.NewMarkAllNotificationsAsReadController(useCase)
}

func BuildDeleteNotificationController(pool *pgxpool.Pool, sender push.Sender) *controllers.DeleteNotificationController {
	repo := adapters.NewPostgreSQLNotificationRepository(pool, sender)
	useCase := application.NewDeleteNotificationUseCase(repo)
	return controllers.NewDeleteNotificationController(useCase)
}

// BuildNotificationPublisher expone el adapter crudo para features que
// necesitan crear notificaciones desde su propio proceso sin pasar por HTTP
// (chat, comentarios, likes) -- lo satisfacen por estructura vía su propia
// interfaz NotificationPublisher (mismo patrón que AssetPhotoProvider).
func BuildNotificationPublisher(pool *pgxpool.Pool, sender push.Sender) *adapters.PostgreSQLNotificationRepository {
	return adapters.NewPostgreSQLNotificationRepository(pool, sender)
}
