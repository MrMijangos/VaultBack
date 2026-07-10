package infrastructure

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"vault/src/features/businesses/application"
	"vault/src/features/businesses/infrastructure/adapters"
	"vault/src/features/businesses/infrastructure/controllers"
)

func BuildCreateBusinessController(pool *pgxpool.Pool) *controllers.CreateBusinessController {
	repo := adapters.NewPostgreSQLBusinessRepository(pool)
	useCase := application.NewCreateBusinessUseCase(repo)
	return controllers.NewCreateBusinessController(useCase)
}

func BuildGetAllBusinessesController(pool *pgxpool.Pool) *controllers.GetAllBusinessesController {
	repo := adapters.NewPostgreSQLBusinessRepository(pool)
	useCase := application.NewGetAllBusinessesUseCase(repo)
	return controllers.NewGetAllBusinessesController(useCase)
}

func BuildGetBusinessByIdController(pool *pgxpool.Pool) *controllers.GetBusinessByIdController {
	repo := adapters.NewPostgreSQLBusinessRepository(pool)
	useCase := application.NewGetBusinessByIdUseCase(repo)
	return controllers.NewGetBusinessByIdController(useCase)
}

func BuildUpdateBusinessController(pool *pgxpool.Pool) *controllers.UpdateBusinessController {
	repo := adapters.NewPostgreSQLBusinessRepository(pool)
	useCase := application.NewUpdateBusinessUseCase(repo)
	return controllers.NewUpdateBusinessController(useCase)
}

func BuildDeleteBusinessController(pool *pgxpool.Pool) *controllers.DeleteBusinessController {
	repo := adapters.NewPostgreSQLBusinessRepository(pool)
	useCase := application.NewDeleteBusinessUseCase(repo)
	return controllers.NewDeleteBusinessController(useCase)
}
