package response

import (
	"time"

	"vault/src/features/businesses/domain/entities"
)

type BusinessPhotoResponse struct {
	ID      string `json:"id"`
	URL     string `json:"url"`
	IsCover bool   `json:"is_cover"`
	Order   int    `json:"order"`
}

type BusinessResponse struct {
	ID           string                  `json:"id"`
	UserID       string                  `json:"user_id"`
	Name         string                  `json:"name"`
	Types        []string                `json:"types"`
	Description  string                  `json:"description"`
	Location     string                  `json:"location"`
	IsVerified   bool                    `json:"is_verified"`
	CreatedAt    time.Time               `json:"created_at"`
	Rating       *float64                `json:"rating"`
	TotalReviews int                     `json:"total_reviews"`
	Specialties  []string                `json:"specialties"`
	Photos       []BusinessPhotoResponse `json:"photos"`
}

func FromEntity(b entities.Business, photos []entities.BusinessPhoto) BusinessResponse {
	photoResponses := make([]BusinessPhotoResponse, 0, len(photos))
	for _, p := range photos {
		photoResponses = append(photoResponses, BusinessPhotoResponse{
			ID:      p.ID,
			URL:     p.URL,
			IsCover: p.IsCover,
			Order:   p.Order,
		})
	}

	return BusinessResponse{
		ID:          b.ID,
		UserID:      b.UserID,
		Name:        b.Name,
		Types:       b.Types,
		Description: b.Description,
		Location:    b.Location,
		IsVerified:  b.IsVerified,
		CreatedAt:   b.CreatedAt,
		Specialties: b.Specialties,
		Photos:      photoResponses,
	}
}
