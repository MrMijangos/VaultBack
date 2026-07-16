package infrastructure

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"vault/src/core/moderation"
	"vault/src/features/reviews/application"
	"vault/src/features/reviews/infrastructure/adapters"
	"vault/src/features/reviews/infrastructure/controllers"
)

func BuildCreateReviewController(pool *pgxpool.Pool, moderationClient *moderation.Client) *controllers.CreateReviewController {
	repo := adapters.NewPostgreSQLReviewRepository(pool)
	useCase := application.NewCreateReviewUseCase(repo, moderationClient)
	return controllers.NewCreateReviewController(useCase)
}

func BuildGetReviewsByProviderController(pool *pgxpool.Pool) *controllers.GetReviewsByProviderController {
	repo := adapters.NewPostgreSQLReviewRepository(pool)
	useCase := application.NewGetReviewsByProviderUseCase(repo)
	return controllers.NewGetReviewsByProviderController(useCase)
}

func BuildGetReviewByIdController(pool *pgxpool.Pool) *controllers.GetReviewByIdController {
	repo := adapters.NewPostgreSQLReviewRepository(pool)
	useCase := application.NewGetReviewByIdUseCase(repo)
	return controllers.NewGetReviewByIdController(useCase)
}

func BuildDeleteReviewController(pool *pgxpool.Pool) *controllers.DeleteReviewController {
	repo := adapters.NewPostgreSQLReviewRepository(pool)
	useCase := application.NewDeleteReviewUseCase(repo)
	return controllers.NewDeleteReviewController(useCase)
}

func BuildLikeReviewController(pool *pgxpool.Pool) *controllers.LikeReviewController {
	repo := adapters.NewPostgreSQLReviewRepository(pool)
	useCase := application.NewLikeReviewUseCase(repo)
	return controllers.NewLikeReviewController(useCase)
}

func BuildUnlikeReviewController(pool *pgxpool.Pool) *controllers.UnlikeReviewController {
	repo := adapters.NewPostgreSQLReviewRepository(pool)
	useCase := application.NewUnlikeReviewUseCase(repo)
	return controllers.NewUnlikeReviewController(useCase)
}
