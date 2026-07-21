package router

import (
	"net/http"

	"vault/src/core/security"
	"vault/src/features/posts/infrastructure/controllers"
)

func RegisterRoutes(
	mux *http.ServeMux,
	createPost *controllers.CreatePostController,
	getAllPosts *controllers.GetAllPostsController,
	getPostById *controllers.GetPostByIdController,
	updatePost *controllers.UpdatePostController,
	deletePost *controllers.DeletePostController,
	uploadPostPhoto *controllers.UploadPostPhotoController,
	likePost *controllers.LikePostController,
	unlikePost *controllers.UnlikePostController,
	savePost *controllers.SavePostController,
	unsavePost *controllers.UnsavePostController,
	getSavedPosts *controllers.GetSavedPostsController,
	jwtSecret string,
) {
	auth := security.RequireAuth(jwtSecret)

	mux.Handle("POST /api/v1/posts", auth(http.HandlerFunc(createPost.Handle)))
	mux.HandleFunc("GET /api/v1/posts", getAllPosts.Handle)
	// Literal más específico que /posts/{id} -- net/http (Go 1.22+) le da
	// prioridad sin importar el orden de registro, pero se deja aquí junto
	// a los otros /posts/{...} por legibilidad.
	mux.Handle("GET /api/v1/posts/saved", auth(http.HandlerFunc(getSavedPosts.Handle)))
	mux.HandleFunc("GET /api/v1/posts/{id}", getPostById.Handle)
	mux.Handle("PUT /api/v1/posts/{id}", auth(http.HandlerFunc(updatePost.Handle)))
	mux.Handle("DELETE /api/v1/posts/{id}", auth(http.HandlerFunc(deletePost.Handle)))
	mux.Handle("POST /api/v1/posts/{id}/photos", auth(http.HandlerFunc(uploadPostPhoto.Handle)))
	mux.Handle("POST /api/v1/posts/{id}/likes", auth(http.HandlerFunc(likePost.Handle)))
	mux.Handle("DELETE /api/v1/posts/{id}/likes", auth(http.HandlerFunc(unlikePost.Handle)))
	mux.Handle("POST /api/v1/posts/{id}/saves", auth(http.HandlerFunc(savePost.Handle)))
	mux.Handle("DELETE /api/v1/posts/{id}/saves", auth(http.HandlerFunc(unsavePost.Handle)))
}
