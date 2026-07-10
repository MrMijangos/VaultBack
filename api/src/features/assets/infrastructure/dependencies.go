package infrastructure

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"vault/src/core/cloudinary"
	"vault/src/features/assets/application"
	"vault/src/features/assets/infrastructure/adapters"
	"vault/src/features/assets/infrastructure/controllers"
)

func BuildCreateAssetController(pool *pgxpool.Pool) *controllers.CreateAssetController {
	repo := adapters.NewPostgreSQLAssetRepository(pool)
	useCase := application.NewCreateAssetUseCase(repo)
	return controllers.NewCreateAssetController(useCase)
}

func BuildGetAllAssetsController(pool *pgxpool.Pool) *controllers.GetAllAssetsController {
	repo := adapters.NewPostgreSQLAssetRepository(pool)
	useCase := application.NewGetAllAssetsUseCase(repo)
	return controllers.NewGetAllAssetsController(useCase)
}

func BuildGetAssetByIdController(pool *pgxpool.Pool) *controllers.GetAssetByIdController {
	repo := adapters.NewPostgreSQLAssetRepository(pool)
	useCase := application.NewGetAssetByIdUseCase(repo)
	return controllers.NewGetAssetByIdController(useCase)
}

func BuildUpdateAssetController(pool *pgxpool.Pool) *controllers.UpdateAssetController {
	repo := adapters.NewPostgreSQLAssetRepository(pool)
	useCase := application.NewUpdateAssetUseCase(repo)
	return controllers.NewUpdateAssetController(useCase)
}

func BuildDeleteAssetController(pool *pgxpool.Pool) *controllers.DeleteAssetController {
	repo := adapters.NewPostgreSQLAssetRepository(pool)
	useCase := application.NewDeleteAssetUseCase(repo)
	return controllers.NewDeleteAssetController(useCase)
}

func BuildUploadAssetPhotoController(pool *pgxpool.Pool, uploader *cloudinary.ImageUploader) *controllers.UploadAssetPhotoController {
	repo := adapters.NewPostgreSQLAssetRepository(pool)
	useCase := application.NewUploadAssetPhotoUseCase(repo, uploader)
	return controllers.NewUploadAssetPhotoController(useCase)
}
