package router

import (
	"net/http"

	"vault/src/core/security"
	"vault/src/features/recommendations/infrastructure/controllers"
)

func RegisterRoutes(
	mux *http.ServeMux,
	getRecommendations *controllers.GetRecommendationsController,
	jwtSecret string,
) {
	auth := security.RequireAuth(jwtSecret)

	mux.Handle("GET /api/v1/recommendations", auth(http.HandlerFunc(getRecommendations.Handle)))
}
