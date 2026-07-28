package repositories

import (
	"context"

	"vault/src/features/recommendations/domain/dto/response"
)

// RecommendationProvider abstrae de dónde vienen las recomendaciones de
// artículos para un usuario. La implementación real (adapters) llama a
// vault-ai-service vía HTTP; esto permite mockear en tests sin necesitar
// el servicio de Python levantado.
type RecommendationProvider interface {
	GetRecommendedItems(ctx context.Context, userID string, topN int) (response.RecommendationsResponse, error)
}
