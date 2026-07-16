package infrastructure

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"vault/src/core/cloudinary"
	"vault/src/features/users/application"
	"vault/src/features/users/infrastructure/adapters"
	"vault/src/features/users/infrastructure/controllers"
)

func BuildCreateUserController(pool *pgxpool.Pool, jwtSecret string, cookieSecure bool) *controllers.CreateUserController {
	repo := adapters.NewPostgreSQLUserRepository(pool)
	useCase := application.NewCreateUserUseCase(repo, jwtSecret)
	return controllers.NewCreateUserController(useCase, cookieSecure)
}

func BuildGetAllUsersController(pool *pgxpool.Pool) *controllers.GetAllUsersController {
	repo := adapters.NewPostgreSQLUserRepository(pool)
	useCase := application.NewGetAllUsersUseCase(repo)
	return controllers.NewGetAllUsersController(useCase)
}

func BuildGetUserByIdController(pool *pgxpool.Pool) *controllers.GetUserByIdController {
	repo := adapters.NewPostgreSQLUserRepository(pool)
	useCase := application.NewGetUserByIdUseCase(repo)
	return controllers.NewGetUserByIdController(useCase)
}

func BuildUpdateUserController(pool *pgxpool.Pool) *controllers.UpdateUserController {
	repo := adapters.NewPostgreSQLUserRepository(pool)
	useCase := application.NewUpdateUserUseCase(repo)
	return controllers.NewUpdateUserController(useCase)
}

func BuildDeleteUserController(pool *pgxpool.Pool) *controllers.DeleteUserController {
	repo := adapters.NewPostgreSQLUserRepository(pool)
	useCase := application.NewDeleteUserUseCase(repo)
	return controllers.NewDeleteUserController(useCase)
}

func BuildUploadUserImageController(pool *pgxpool.Pool, uploader *cloudinary.ImageUploader) *controllers.UploadUserImageController {
	repo := adapters.NewPostgreSQLUserRepository(pool)
	useCase := application.NewUploadUserImageUseCase(repo, uploader)
	return controllers.NewUploadUserImageController(useCase)
}
