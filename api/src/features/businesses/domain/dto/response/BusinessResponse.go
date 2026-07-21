package response

import (
	"time"

	"vault/src/features/businesses/domain/entities"
)

type BusinessResponse struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	Name         string    `json:"name"`
	Type         string    `json:"type"`
	Description  string    `json:"description"`
	Location     string    `json:"location"`
	IsVerified   bool      `json:"is_verified"`
	CreatedAt    time.Time `json:"created_at"`
	Rating       *float64  `json:"rating"`
	TotalReviews int       `json:"total_reviews"`
}

func FromEntity(b entities.Business) BusinessResponse {
	return BusinessResponse{
		ID:          b.ID,
		UserID:      b.UserID,
		Name:        b.Name,
		Type:        b.Type,
		Description: b.Description,
		Location:    b.Location,
		IsVerified:  b.IsVerified,
		CreatedAt:   b.CreatedAt,
	}
}

func FromEntities(list []entities.Business) []BusinessResponse {
	out := make([]BusinessResponse, 0, len(list))
	for _, b := range list {
		out = append(out, FromEntity(b))
	}
	return out
}
