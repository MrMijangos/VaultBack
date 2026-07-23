package router

import (
	"net/http"

	"vault/src/core/security"
	"vault/src/features/addresses/infrastructure/controllers"
)

func RegisterRoutes(
	mux *http.ServeMux,
	createAddress *controllers.CreateAddressController,
	listAddresses *controllers.ListAddressesController,
	deleteAddress *controllers.DeleteAddressController,
	setDefaultAddress *controllers.SetDefaultAddressController,
	jwtSecret string,
) {
	auth := security.RequireAuth(jwtSecret)

	mux.Handle("GET /api/v1/addresses", auth(http.HandlerFunc(listAddresses.Handle)))
	mux.Handle("POST /api/v1/addresses", auth(http.HandlerFunc(createAddress.Handle)))
	mux.Handle("DELETE /api/v1/addresses/{id}", auth(http.HandlerFunc(deleteAddress.Handle)))
	mux.Handle("PATCH /api/v1/addresses/{id}/default", auth(http.HandlerFunc(setDefaultAddress.Handle)))
}
