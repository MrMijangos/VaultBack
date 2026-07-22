package infrastructure

import (
	"github.com/jackc/pgx/v5/pgxpool"

	businessesAdapters "vault/src/features/businesses/infrastructure/adapters"

	"vault/src/features/businessservices/application"
	"vault/src/features/businessservices/infrastructure/adapters"
	"vault/src/features/businessservices/infrastructure/controllers"
)

func BuildCreateBusinessServiceController(pool *pgxpool.Pool) *controllers.CreateBusinessServiceController {
	repo := adapters.NewPostgreSQLBusinessServiceRepository(pool)
	ownerProvider := adapters.NewBusinessOwnerAdapter(businessesAdapters.NewPostgreSQLBusinessRepository(pool))
	useCase := application.NewCreateBusinessServiceUseCase(repo, ownerProvider)
	return controllers.NewCreateBusinessServiceController(useCase)
}

func BuildListBusinessServicesController(pool *pgxpool.Pool) *controllers.ListBusinessServicesController {
	repo := adapters.NewPostgreSQLBusinessServiceRepository(pool)
	useCase := application.NewListBusinessServicesUseCase(repo)
	return controllers.NewListBusinessServicesController(useCase)
}

func BuildUpdateBusinessServiceController(pool *pgxpool.Pool) *controllers.UpdateBusinessServiceController {
	repo := adapters.NewPostgreSQLBusinessServiceRepository(pool)
	ownerProvider := adapters.NewBusinessOwnerAdapter(businessesAdapters.NewPostgreSQLBusinessRepository(pool))
	useCase := application.NewUpdateBusinessServiceUseCase(repo, ownerProvider)
	return controllers.NewUpdateBusinessServiceController(useCase)
}

func BuildDeleteBusinessServiceController(pool *pgxpool.Pool) *controllers.DeleteBusinessServiceController {
	repo := adapters.NewPostgreSQLBusinessServiceRepository(pool)
	ownerProvider := adapters.NewBusinessOwnerAdapter(businessesAdapters.NewPostgreSQLBusinessRepository(pool))
	useCase := application.NewDeleteBusinessServiceUseCase(repo, ownerProvider)
	return controllers.NewDeleteBusinessServiceController(useCase)
}
