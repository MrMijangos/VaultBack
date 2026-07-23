package router

import (
	"net/http"

	"vault/src/core/security"
	"vault/src/features/assetcomments/infrastructure/controllers"
)

func RegisterRoutes(
	mux *http.ServeMux,
	createAssetComment *controllers.CreateAssetCommentController,
	getAssetComments *controllers.GetAssetCommentsController,
	deleteAssetComment *controllers.DeleteAssetCommentController,
	jwtSecret string,
) {
	auth := security.RequireAuth(jwtSecret)

	mux.Handle("POST /api/v1/assets/{id}/comments", auth(http.HandlerFunc(createAssetComment.Handle)))
	mux.HandleFunc("GET /api/v1/assets/{id}/comments", getAssetComments.Handle)
	mux.Handle("DELETE /api/v1/assets/{id}/comments/{commentId}", auth(http.HandlerFunc(deleteAssetComment.Handle)))
}
