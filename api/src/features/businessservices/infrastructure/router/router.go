package router

import (
	"net/http"

	"vault/src/core/security"
	"vault/src/features/businessservices/infrastructure/controllers"
)

func RegisterRoutes(
	mux *http.ServeMux,
	createBusinessService *controllers.CreateBusinessServiceController,
	listBusinessServices *controllers.ListBusinessServicesController,
	updateBusinessService *controllers.UpdateBusinessServiceController,
	deleteBusinessService *controllers.DeleteBusinessServiceController,
	jwtSecret string,
) {
	auth := security.RequireAuth(jwtSecret)

	mux.Handle("POST /api/v1/businesses/{id}/services", auth(http.HandlerFunc(createBusinessService.Handle)))
	mux.HandleFunc("GET /api/v1/businesses/{id}/services", listBusinessServices.Handle)
	mux.Handle("PUT /api/v1/businesses/{id}/services/{serviceId}", auth(http.HandlerFunc(updateBusinessService.Handle)))
	mux.Handle("DELETE /api/v1/businesses/{id}/services/{serviceId}", auth(http.HandlerFunc(deleteBusinessService.Handle)))
}
