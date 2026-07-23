package response

import (
	"time"

	"vault/src/features/assetcomments/domain/entities"
)

type AssetCommentResponse struct {
	ID              string    `json:"id"`
	AssetID         string    `json:"asset_id"`
	UserID          string    `json:"user_id"`
	Content         string    `json:"content"`
	IsVisible       bool      `json:"is_visible"`
	CreatedAt       time.Time `json:"created_at"`
	AuthorName      string    `json:"author_name"`
	AuthorAvatarURL string    `json:"author_avatar_url"`
}

func FromEntity(c entities.AssetComment) AssetCommentResponse {
	return AssetCommentResponse{
		ID:              c.ID,
		AssetID:         c.AssetID,
		UserID:          c.UserID,
		Content:         c.Content,
		IsVisible:       c.IsVisible,
		CreatedAt:       c.CreatedAt,
		AuthorName:      c.AuthorName,
		AuthorAvatarURL: c.AuthorAvatarURL,
	}
}

func FromEntities(list []entities.AssetComment) []AssetCommentResponse {
	out := make([]AssetCommentResponse, 0, len(list))
	for _, c := range list {
		out = append(out, FromEntity(c))
	}
	return out
}
