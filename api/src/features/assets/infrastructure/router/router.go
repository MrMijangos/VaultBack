package router

import (
	"net/http"

	"vault/src/core/security"
	"vault/src/features/assets/infrastructure/controllers"
)

func RegisterRoutes(
	mux *http.ServeMux,
	createAsset *controllers.CreateAssetController,
	getAllAssets *controllers.GetAllAssetsController,
	getAssetById *controllers.GetAssetByIdController,
	updateAsset *controllers.UpdateAssetController,
	deleteAsset *controllers.DeleteAssetController,
	uploadAssetPhoto *controllers.UploadAssetPhotoController,
	jwtSecret string,
) {
	auth := security.RequireAuth(jwtSecret)

	mux.Handle("POST /api/v1/assets", auth(http.HandlerFunc(createAsset.Handle)))
	mux.HandleFunc("GET /api/v1/assets", getAllAssets.Handle)
	mux.HandleFunc("GET /api/v1/assets/{id}", getAssetById.Handle)
	mux.Handle("PUT /api/v1/assets/{id}", auth(http.HandlerFunc(updateAsset.Handle)))
	mux.Handle("DELETE /api/v1/assets/{id}", auth(http.HandlerFunc(deleteAsset.Handle)))
	mux.Handle("POST /api/v1/assets/{id}/photos", auth(http.HandlerFunc(uploadAssetPhoto.Handle)))
}
