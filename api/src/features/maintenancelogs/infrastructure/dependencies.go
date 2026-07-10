package infrastructure

import (
	"github.com/jackc/pgx/v5/pgxpool"

	assetadapters "vault/src/features/assets/infrastructure/adapters"
	"vault/src/features/maintenancelogs/application"
	"vault/src/features/maintenancelogs/infrastructure/adapters"
	"vault/src/features/maintenancelogs/infrastructure/controllers"
)

func BuildCreateMaintenanceLogController(pool *pgxpool.Pool) *controllers.CreateMaintenanceLogController {
	repo := adapters.NewPostgreSQLMaintenanceLogRepository(pool)
	assetRepo := assetadapters.NewPostgreSQLAssetRepository(pool)
	useCase := application.NewCreateMaintenanceLogUseCase(repo, assetRepo)
	return controllers.NewCreateMaintenanceLogController(useCase)
}

func BuildGetLogsByAssetController(pool *pgxpool.Pool) *controllers.GetLogsByAssetController {
	repo := adapters.NewPostgreSQLMaintenanceLogRepository(pool)
	useCase := application.NewGetLogsByAssetUseCase(repo)
	return controllers.NewGetLogsByAssetController(useCase)
}

func BuildGetMaintenanceLogByIdController(pool *pgxpool.Pool) *controllers.GetMaintenanceLogByIdController {
	repo := adapters.NewPostgreSQLMaintenanceLogRepository(pool)
	useCase := application.NewGetMaintenanceLogByIdUseCase(repo)
	return controllers.NewGetMaintenanceLogByIdController(useCase)
}

func BuildUpdateMaintenanceLogController(pool *pgxpool.Pool) *controllers.UpdateMaintenanceLogController {
	repo := adapters.NewPostgreSQLMaintenanceLogRepository(pool)
	assetRepo := assetadapters.NewPostgreSQLAssetRepository(pool)
	useCase := application.NewUpdateMaintenanceLogUseCase(repo, assetRepo)
	return controllers.NewUpdateMaintenanceLogController(useCase)
}

func BuildDeleteMaintenanceLogController(pool *pgxpool.Pool) *controllers.DeleteMaintenanceLogController {
	repo := adapters.NewPostgreSQLMaintenanceLogRepository(pool)
	assetRepo := assetadapters.NewPostgreSQLAssetRepository(pool)
	useCase := application.NewDeleteMaintenanceLogUseCase(repo, assetRepo)
	return controllers.NewDeleteMaintenanceLogController(useCase)
}
