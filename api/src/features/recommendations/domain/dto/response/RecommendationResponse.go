package response

// RecommendedItemResponse es un artículo recomendado del catálogo.
// Los campos coinciden con RecommendedItemDTO de vault-ai-service
// (src/domain/dtos/recommended_item_dto.py).
type RecommendedItemResponse struct {
	ItemID   string  `json:"item_id"`
	Category string  `json:"category"`
	Brand    string  `json:"brand"`
	Model    string  `json:"model"`
	Score    float64 `json:"score"`
	Strategy string  `json:"strategy"` // "content" | "cluster"
}

// RecommendationsResponse envuelve la lista de items recomendados para un
// usuario. Coincide con RecommendItemsResponseDTO del servicio de ML.
type RecommendationsResponse struct {
	UserID string                    `json:"user_id"`
	Items  []RecommendedItemResponse `json:"items"`
}
