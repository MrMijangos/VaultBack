package response

import (
	"time"

	"vault-payment/src/features/ads/domain/entities"
)

type AdResponse struct {
	ID            string    `json:"id"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	ImageURL      string    `json:"image_url"`
	TargetSection string    `json:"target_section"`
	TargetID      string    `json:"target_id"`
	Status        string    `json:"status"`
	Impressions   int64     `json:"impressions"`
	Clicks        int64     `json:"clicks"`
	CreatedAt     time.Time `json:"created_at"`
}

func AdFromEntity(a *entities.Ad) AdResponse {
	return AdResponse{
		ID:            a.ID,
		Title:         a.Title,
		Description:   a.Description,
		ImageURL:      a.ImageURL,
		TargetSection: a.TargetSection,
		TargetID:      a.TargetID,
		Status:        a.Status,
		Impressions:   a.Impressions,
		Clicks:        a.Clicks,
		CreatedAt:     a.CreatedAt,
	}
}

func AdsFromEntities(ads []*entities.Ad) []AdResponse {
	out := make([]AdResponse, 0, len(ads))
	for _, a := range ads {
		out = append(out, AdFromEntity(a))
	}
	return out
}
