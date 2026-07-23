package infrastructure

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"vault/src/core/moderation"
	"vault/src/features/assetcomments/application"
	"vault/src/features/assetcomments/infrastructure/adapters"
	"vault/src/features/assetcomments/infrastructure/controllers"
)

func BuildCreateAssetCommentController(pool *pgxpool.Pool, moderationClient *moderation.Client) *controllers.CreateAssetCommentController {
	repo := adapters.NewPostgreSQLAssetCommentRepository(pool)
	useCase := application.NewCreateAssetCommentUseCase(repo, moderationClient)
	return controllers.NewCreateAssetCommentController(useCase)
}

func BuildGetAssetCommentsController(pool *pgxpool.Pool) *controllers.GetAssetCommentsController {
	repo := adapters.NewPostgreSQLAssetCommentRepository(pool)
	useCase := application.NewGetAssetCommentsUseCase(repo)
	return controllers.NewGetAssetCommentsController(useCase)
}

func BuildDeleteAssetCommentController(pool *pgxpool.Pool) *controllers.DeleteAssetCommentController {
	repo := adapters.NewPostgreSQLAssetCommentRepository(pool)
	useCase := application.NewDeleteAssetCommentUseCase(repo)
	return controllers.NewDeleteAssetCommentController(useCase)
}
