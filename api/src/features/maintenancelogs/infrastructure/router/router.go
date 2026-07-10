package router

import (
	"net/http"

	"vault/src/core/security"
	"vault/src/features/maintenancelogs/infrastructure/controllers"
)

func RegisterRoutes(
	mux *http.ServeMux,
	createLog *controllers.CreateMaintenanceLogController,
	getLogsByAsset *controllers.GetLogsByAssetController,
	getLogById *controllers.GetMaintenanceLogByIdController,
	updateLog *controllers.UpdateMaintenanceLogController,
	deleteLog *controllers.DeleteMaintenanceLogController,
	jwtSecret string,
) {
	auth := security.RequireAuth(jwtSecret)

	mux.Handle("POST /api/v1/maintenance-logs", auth(http.HandlerFunc(createLog.Handle)))
	mux.HandleFunc("GET /api/v1/maintenance-logs", getLogsByAsset.Handle)
	mux.HandleFunc("GET /api/v1/maintenance-logs/{id}", getLogById.Handle)
	mux.Handle("PUT /api/v1/maintenance-logs/{id}", auth(http.HandlerFunc(updateLog.Handle)))
	mux.Handle("DELETE /api/v1/maintenance-logs/{id}", auth(http.HandlerFunc(deleteLog.Handle)))
}
