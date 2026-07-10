package infrastructure

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"vault/src/features/comments/application"
	"vault/src/features/comments/infrastructure/adapters"
	"vault/src/features/comments/infrastructure/controllers"
)

func BuildCreateCommentController(pool *pgxpool.Pool) *controllers.CreateCommentController {
	repo := adapters.NewPostgreSQLCommentRepository(pool)
	useCase := application.NewCreateCommentUseCase(repo)
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
