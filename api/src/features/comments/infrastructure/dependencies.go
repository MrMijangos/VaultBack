package infrastructure

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"vault/src/core/moderation"
	"vault/src/core/push"
	"vault/src/features/comments/application"
	"vault/src/features/comments/infrastructure/adapters"
	"vault/src/features/comments/infrastructure/controllers"
	notificationsInfra "vault/src/features/notifications/infrastructure"
	postsAdapters "vault/src/features/posts/infrastructure/adapters"
)

func BuildCreateCommentController(pool *pgxpool.Pool, moderationClient *moderation.Client, sender push.Sender) *controllers.CreateCommentController {
	repo := adapters.NewPostgreSQLCommentRepository(pool)
	postAuthors := postsAdapters.NewPostgreSQLPostRepository(pool)
	notifier := notificationsInfra.BuildNotificationPublisher(pool, sender)
	useCase := application.NewCreateCommentUseCase(repo, moderationClient, postAuthors, notifier)
	return controllers.NewCreateCommentController(useCase)
}

func BuildGetCommentsByPostController(pool *pgxpool.Pool) *controllers.GetCommentsByPostController {
	repo := adapters.NewPostgreSQLCommentRepository(pool)
	useCase := application.NewGetCommentsByPostUseCase(repo)
	return controllers.NewGetCommentsByPostController(useCase)
}

func BuildDeleteCommentController(pool *pgxpool.Pool) *controllers.DeleteCommentController {
	repo := adapters.NewPostgreSQLCommentRepository(pool)
	useCase := application.NewDeleteCommentUseCase(repo)
	return controllers.NewDeleteCommentController(useCase)
}
