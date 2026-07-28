package infrastructure

import (
	"vault/src/features/recommendations/application"
	"vault/src/features/recommendations/infrastructure/adapters"
	"vault/src/features/recommendations/infrastructure/controllers"
)

// BuildGetRecommendationsController arma el controller con el cliente HTTP
// hacia vault-ai-service. aiServiceURL debe ser la misma URL que ya usa
// cfg.NLPServiceURL (main.go) -- es el mismo servicio, solo otro endpoint.
func BuildGetRecommendationsController(aiServiceURL string) *controllers.GetRecommendationsController {
	provider := adapters.NewHTTPRecommendationClient(aiServiceURL)
	useCase := application.NewGetRecommendationsUseCase(provider)
	return controllers.NewGetRecommendationsController(useCase)
}
