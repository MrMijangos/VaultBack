package infrastructure

import (
	"github.com/jackc/pgx/v5/pgxpool"

	reviewsAdapters "vault/src/features/reviews/infrastructure/adapters"

	"vault/src/features/restorerprofiles/application"
	"vault/src/features/restorerprofiles/infrastructure/adapters"
	"vault/src/features/restorerprofiles/infrastructure/controllers"
)

func BuildUpsertRestorerProfileController(pool *pgxpool.Pool) *controllers.UpsertRestorerProfileController {
	repo := adapters.NewPostgreSQLRestorerProfileRepository(pool)
	useCase := application.NewUpsertRestorerProfileUseCase(repo)
	return controllers.NewUpsertRestorerProfileController(useCase)
}

func BuildGetRestorerProfileController(pool *pgxpool.Pool) *controllers.GetRestorerProfileController {
	repo := adapters.NewPostgreSQLRestorerProfileRepository(pool)
	ratingRepo := reviewsAdapters.NewPostgreSQLReviewRepository(pool)
	useCase := application.NewGetRestorerProfileUseCase(repo, ratingRepo)
	return controllers.NewGetRestorerProfileController(useCase)
}

func BuildListRestorerProfilesController(pool *pgxpool.Pool) *controllers.ListRestorerProfilesController {
	repo := adapters.NewPostgreSQLRestorerProfileRepository(pool)
	ratingRepo := reviewsAdapters.NewPostgreSQLReviewRepository(pool)
	useCase := application.NewListRestorerProfilesUseCase(repo, ratingRepo)
	return controllers.NewListRestorerProfilesController(useCase)
}
