package response

import (
	"vault/src/features/restorerprofiles/domain/entities"
)

type RestorerServiceResponse struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

type RestorerProfileResponse struct {
	UserID       string                    `json:"user_id"`
	Name         string                    `json:"name"`
	AvatarURL    string                    `json:"avatar_url"`
	Bio          string                    `json:"bio"`
	Specialties  []string                  `json:"specialties"`
	Services     []RestorerServiceResponse `json:"services"`
	Rating       *float64                  `json:"rating"`
	ReviewsCount int                       `json:"reviews_count"`
}

func FromEntity(p entities.RestorerProfile) RestorerProfileResponse {
	services := make([]RestorerServiceResponse, 0, len(p.Services))
	for _, s := range p.Services {
		services = append(services, RestorerServiceResponse{
			ID:          s.ID,
			Title:       s.Title,
			Description: s.Description,
			Price:       s.Price,
		})
	}

	specialties := p.Specialties
	if specialties == nil {
		specialties = []string{}
	}

	return RestorerProfileResponse{
		UserID:       p.UserID,
		Name:         p.Name,
		AvatarURL:    p.AvatarURL,
		Bio:          p.Bio,
		Specialties:  specialties,
		Services:     services,
		Rating:       p.Rating,
		ReviewsCount: p.TotalReviews,
	}
}

func FromEntities(list []entities.RestorerProfile) []RestorerProfileResponse {
	out := make([]RestorerProfileResponse, 0, len(list))
	for _, p := range list {
		out = append(out, FromEntity(p))
	}
	return out
}
