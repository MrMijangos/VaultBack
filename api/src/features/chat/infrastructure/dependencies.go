package infrastructure

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"vault/src/core/push"
	"vault/src/features/chat/application"
	"vault/src/features/chat/infrastructure/adapters"
	"vault/src/features/chat/infrastructure/controllers"
	notificationsInfra "vault/src/features/notifications/infrastructure"
)

// BuildSendChatMessageController no recibe un moderationClient -- a
// diferencia de comments/reviews/posts, el chat es E2EE y el servidor nunca
// ve el contenido en texto plano, por lo que no hay nada que moderar.
func BuildSendChatMessageController(pool *pgxpool.Pool, sender push.Sender) *controllers.SendChatMessageController {
	repo := adapters.NewPostgreSQLChatMessageRepository(pool)
	notifier := notificationsInfra.BuildNotificationPublisher(pool, sender)
	useCase := application.NewSendChatMessageUseCase(repo, notifier)
	return controllers.NewSendChatMessageController(useCase)
}

func BuildGetConversationMessagesController(pool *pgxpool.Pool) *controllers.GetConversationMessagesController {
	repo := adapters.NewPostgreSQLChatMessageRepository(pool)
	useCase := application.NewGetConversationMessagesUseCase(repo)
	return controllers.NewGetConversationMessagesController(useCase)
}

func BuildUpdateChatMessageStatusController(pool *pgxpool.Pool) *controllers.UpdateChatMessageStatusController {
	repo := adapters.NewPostgreSQLChatMessageRepository(pool)
	useCase := application.NewUpdateChatMessageStatusUseCase(repo)
	return controllers.NewUpdateChatMessageStatusController(useCase)
}

func BuildGetConversationsController(pool *pgxpool.Pool) *controllers.GetConversationsController {
	repo := adapters.NewPostgreSQLChatMessageRepository(pool)
	useCase := application.NewGetConversationsUseCase(repo)
	return controllers.NewGetConversationsController(useCase)
}

func BuildDeleteChatMessageController(pool *pgxpool.Pool) *controllers.DeleteChatMessageController {
	repo := adapters.NewPostgreSQLChatMessageRepository(pool)
	useCase := application.NewDeleteChatMessageUseCase(repo)
	return controllers.NewDeleteChatMessageController(useCase)
}

func BuildDeleteConversationController(pool *pgxpool.Pool) *controllers.DeleteConversationController {
	repo := adapters.NewPostgreSQLChatMessageRepository(pool)
	useCase := application.NewDeleteConversationUseCase(repo)
	return controllers.NewDeleteConversationController(useCase)
}
