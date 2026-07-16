package response

import (
	"time"

	"vault/src/features/reviews/domain/entities"
)

type ReviewResponse struct {
	ID              string    `json:"id"`
	UserID          string    `json:"user_id"`
	ProviderID      string    `json:"provider_id"`
	Content         string    `json:"content"`
	IsVisible       bool      `json:"is_visible"`
	LikesCount      int       `json:"likes_count"`
	CreatedAt       time.Time `json:"created_at"`
	AuthorName      string    `json:"author_name"`
	AuthorAvatarURL string    `json:"author_avatar_url"`
}

func FromEntity(r entities.Review) ReviewResponse {
	return ReviewResponse{
		ID:              r.ID,
		UserID:          r.UserID,
		ProviderID:      r.ProviderID,
		Content:         r.Content,
		IsVisible:       r.IsVisible,
		LikesCount:      r.LikesCount,
		CreatedAt:       r.CreatedAt,
		AuthorName:      r.AuthorName,
		AuthorAvatarURL: r.AuthorAvatarURL,
	}
}

func FromEntities(list []entities.Review) []ReviewResponse {
	out := make([]ReviewResponse, 0, len(list))
	for _, r := range list {
		out = append(out, FromEntity(r))
	}
	return out
}
