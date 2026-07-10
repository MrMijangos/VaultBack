package router

import (
	"net/http"

	"vault/src/core/security"
	"vault/src/features/comments/infrastructure/controllers"
)

func RegisterRoutes(
	mux *http.ServeMux,
	createComment *controllers.CreateCommentController,
	getCommentsByPost *controllers.GetCommentsByPostController,
	deleteComment *controllers.DeleteCommentController,
	jwtSecret string,
) {
	auth := security.RequireAuth(jwtSecret)

	mux.Handle("POST /api/v1/posts/{id}/comments", auth(http.HandlerFunc(createComment.Handle)))
	mux.HandleFunc("GET /api/v1/posts/{id}/comments", getCommentsByPost.Handle)
	mux.Handle("DELETE /api/v1/comments/{id}", auth(http.HandlerFunc(deleteComment.Handle)))
}
