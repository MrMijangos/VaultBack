package response

import (
	"time"

	"vault/src/features/businessservices/domain/entities"
)

type BusinessServiceResponse struct {
	ID          string    `json:"id"`
	BusinessID  string    `json:"business_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	CreatedAt   time.Time `json:"created_at"`
}

func FromEntity(s entities.BusinessService) BusinessServiceResponse {
	return BusinessServiceResponse{
		ID:          s.ID,
		BusinessID:  s.BusinessID,
		Title:       s.Title,
		Description: s.Description,
		Price:       s.Price,
		CreatedAt:   s.CreatedAt,
	}
}

func FromEntities(list []entities.BusinessService) []BusinessServiceResponse {
	out := make([]BusinessServiceResponse, 0, len(list))
	for _, s := range list {
		out = append(out, FromEntity(s))
	}
	return out
}
