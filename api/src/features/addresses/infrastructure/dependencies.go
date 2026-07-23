package infrastructure

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"vault/src/features/addresses/application"
	"vault/src/features/addresses/infrastructure/adapters"
	"vault/src/features/addresses/infrastructure/controllers"
)

func BuildCreateAddressController(pool *pgxpool.Pool) *controllers.CreateAddressController {
	repo := adapters.NewPostgreSQLAddressRepository(pool)
	useCase := application.NewCreateAddressUseCase(repo)
	return controllers.NewCreateAddressController(useCase)
}

func BuildListAddressesController(pool *pgxpool.Pool) *controllers.ListAddressesController {
	repo := adapters.NewPostgreSQLAddressRepository(pool)
	useCase := application.NewListAddressesUseCase(repo)
	return controllers.NewListAddressesController(useCase)
}

func BuildDeleteAddressController(pool *pgxpool.Pool) *controllers.DeleteAddressController {
	repo := adapters.NewPostgreSQLAddressRepository(pool)
	useCase := application.NewDeleteAddressUseCase(repo)
	return controllers.NewDeleteAddressController(useCase)
}

func BuildSetDefaultAddressController(pool *pgxpool.Pool) *controllers.SetDefaultAddressController {
	repo := adapters.NewPostgreSQLAddressRepository(pool)
	useCase := application.NewSetDefaultAddressUseCase(repo)
	return controllers.NewSetDefaultAddressController(useCase)
}
