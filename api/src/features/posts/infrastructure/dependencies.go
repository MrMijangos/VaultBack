package infrastructure

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"vault/src/core/cloudinary"
	"vault/src/features/posts/application"
	"vault/src/features/posts/infrastructure/adapters"
	"vault/src/features/posts/infrastructure/controllers"
)

func BuildCreatePostController(pool *pgxpool.Pool) *controllers.CreatePostController {
	repo := adapters.NewPostgreSQLPostRepository(pool)
	useCase := application.NewCreatePostUseCase(repo)
	return controllers.NewCreatePostController(useCase)
}

func BuildGetAllPostsController(pool *pgxpool.Pool) *controllers.GetAllPostsController {
	repo := adapters.NewPostgreSQLPostRepository(pool)
	useCase := application.NewGetAllPostsUseCase(repo)
	return controllers.NewGetAllPostsController(useCase)
}

func BuildGetPostByIdController(pool *pgxpool.Pool) *controllers.GetPostByIdController {
	repo := adapters.NewPostgreSQLPostRepository(pool)
	useCase := application.NewGetPostByIdUseCase(repo)
	return controllers.NewGetPostByIdController(useCase)
}

func BuildUpdatePostController(pool *pgxpool.Pool) *controllers.UpdatePostController {
	repo := adapters.NewPostgreSQLPostRepository(pool)
	useCase := application.NewUpdatePostUseCase(repo)
	return controllers.NewUpdatePostController(useCase)
}

func BuildDeletePostController(pool *pgxpool.Pool) *controllers.DeletePostController {
	repo := adapters.NewPostgreSQLPostRepository(pool)
	useCase := application.NewDeletePostUseCase(repo)
	return controllers.NewDeletePostController(useCase)
}

func BuildUploadPostPhotoController(pool *pgxpool.Pool, uploader *cloudinary.ImageUploader) *controllers.UploadPostPhotoController {
	repo := adapters.NewPostgreSQLPostRepository(pool)
	useCase := application.NewUploadPostPhotoUseCase(repo, uploader)
	return controllers.NewUploadPostPhotoController(useCase)
}

func BuildLikePostController(pool *pgxpool.Pool) *controllers.LikePostController {
	repo := adapters.NewPostgreSQLPostRepository(pool)
	useCase := application.NewLikePostUseCase(repo)
	return controllers.NewLikePostController(useCase)
}

func BuildUnlikePostController(pool *pgxpool.Pool) *controllers.UnlikePostController {
	repo := adapters.NewPostgreSQLPostRepository(pool)
	useCase := application.NewUnlikePostUseCase(repo)
	return controllers.NewUnlikePostController(useCase)
}
