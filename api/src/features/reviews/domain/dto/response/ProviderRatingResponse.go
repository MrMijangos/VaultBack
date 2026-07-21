package response

type ProviderRatingResponse struct {
	ProviderID   string   `json:"provider_id"`
	Rating       *float64 `json:"rating"`
	TotalReviews int      `json:"total_reviews"`
}
