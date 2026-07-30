package router

import (
	"net/http"

	"vault/src/core/security"
	"vault/src/features/servicerequests/infrastructure/controllers"
)

func RegisterRoutes(
	mux *http.ServeMux,
	create *controllers.CreateServiceRequestController,
	accept *controllers.AcceptServiceRequestController,
	start *controllers.StartServiceRequestController,
	finish *controllers.FinishServiceRequestController,
	confirm *controllers.ConfirmServiceRequestController,
	listMine *controllers.ListMyServiceRequestsController,
	listIncoming *controllers.ListIncomingServiceRequestsController,
	jwtSecret string,
) {
	auth := security.RequireAuth(jwtSecret)

	mux.Handle("POST /api/v1/service-requests", auth(http.HandlerFunc(create.Handle)))
	mux.Handle("GET /api/v1/service-requests/mine", auth(http.HandlerFunc(listMine.Handle)))
	mux.Handle("GET /api/v1/service-requests/incoming", auth(http.HandlerFunc(listIncoming.Handle)))
	mux.Handle("POST /api/v1/service-requests/{id}/accept", auth(http.HandlerFunc(accept.Handle)))
	mux.Handle("POST /api/v1/service-requests/{id}/start", auth(http.HandlerFunc(start.Handle)))
	mux.Handle("POST /api/v1/service-requests/{id}/finish", auth(http.HandlerFunc(finish.Handle)))
	mux.Handle("POST /api/v1/service-requests/{id}/confirm", auth(http.HandlerFunc(confirm.Handle)))
}
