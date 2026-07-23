package router

import (
	"net/http"

	"vault/src/core/security"
	"vault/src/features/restorerprofiles/infrastructure/controllers"
)

func RegisterRoutes(
	mux *http.ServeMux,
	upsertRestorerProfile *controllers.UpsertRestorerProfileController,
	getRestorerProfile *controllers.GetRestorerProfileController,
	listRestorerProfiles *controllers.ListRestorerProfilesController,
	jwtSecret string,
) {
	auth := security.RequireAuth(jwtSecret)

	mux.HandleFunc("GET /api/v1/restorerprofiles", listRestorerProfiles.Handle)
	mux.HandleFunc("GET /api/v1/restorerprofiles/{userId}", getRestorerProfile.Handle)
	mux.Handle("PUT /api/v1/restorerprofiles/{userId}", auth(http.HandlerFunc(upsertRestorerProfile.Handle)))
}
