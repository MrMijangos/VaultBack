package infrastructure

import (
	"github.com/jackc/pgx/v5/pgxpool"

	assetadapters "vault/src/features/assets/infrastructure/adapters"
	"vault/src/features/blockchaincertificates/application"
	"vault/src/features/blockchaincertificates/infrastructure/adapters"
	"vault/src/features/blockchaincertificates/infrastructure/controllers"
)

func BuildCreateBlockchainCertificateController(pool *pgxpool.Pool) *controllers.CreateBlockchainCertificateController {
	repo := adapters.NewPostgreSQLBlockchainCertificateRepository(pool)
	assetRepo := assetadapters.NewPostgreSQLAssetRepository(pool)
	useCase := application.NewCreateBlockchainCertificateUseCase(repo, assetRepo)
	return controllers.NewCreateBlockchainCertificateController(useCase)
}

func BuildGetCertificatesByAssetController(pool *pgxpool.Pool) *controllers.GetCertificatesByAssetController {
	repo := adapters.NewPostgreSQLBlockchainCertificateRepository(pool)
	useCase := application.NewGetCertificatesByAssetUseCase(repo)
	return controllers.NewGetCertificatesByAssetController(useCase)
}

func BuildGetCertificateByIdController(pool *pgxpool.Pool) *controllers.GetCertificateByIdController {
	repo := adapters.NewPostgreSQLBlockchainCertificateRepository(pool)
	useCase := application.NewGetCertificateByIdUseCase(repo)
	return controllers.NewGetCertificateByIdController(useCase)
}
