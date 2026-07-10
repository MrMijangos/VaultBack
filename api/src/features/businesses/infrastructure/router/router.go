package router

import (
	"net/http"

	"vault/src/core/security"
	"vault/src/features/businesses/infrastructure/controllers"
)

func RegisterRoutes(
	mux *http.ServeMux,
	createBusiness *controllers.CreateBusinessController,
	getAllBusinesses *controllers.GetAllBusinessesController,
	getBusinessById *controllers.GetBusinessByIdController,
	updateBusiness *controllers.UpdateBusinessController,
	deleteBusiness *controllers.DeleteBusinessController,
	jwtSecret string,
) {
	auth := security.RequireAuth(jwtSecret)

	mux.Handle("POST /api/v1/businesses", auth(http.HandlerFunc(createBusiness.Handle)))
	mux.HandleFunc("GET /api/v1/businesses", getAllBusinesses.Handle)
	mux.HandleFunc("GET /api/v1/businesses/{id}", getBusinessById.Handle)
	mux.Handle("PUT /api/v1/businesses/{id}", auth(http.HandlerFunc(updateBusiness.Handle)))
	mux.Handle("DELETE /api/v1/businesses/{id}", auth(http.HandlerFunc(deleteBusiness.Handle)))
}
