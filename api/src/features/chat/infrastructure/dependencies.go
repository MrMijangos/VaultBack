package infrastructure

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"vault/src/features/chat/application"
	"vault/src/features/chat/infrastructure/adapters"
	"vault/src/features/chat/infrastructure/controllers"
)

// BuildSendChatMessageController no recibe un moderationClient -- a
// diferencia de comments/reviews/posts, el chat es E2EE y el servidor nunca
// ve el contenido en texto plano, por lo que no hay nada que moderar.
func BuildSendChatMessageController(pool *pgxpool.Pool) *controllers.SendChatMessageController {
	repo := adapters.NewPostgreSQLChatMessageRepository(pool)
	useCase := application.NewSendChatMessageUseCase(repo)
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
