package response

import (
	"time"

	"vault/src/features/posts/domain/entities"
)

type PostPhotoResponse struct {
	ID    string `json:"id"`
	URL   string `json:"url"`
	Order int    `json:"order"`
}

type PostResponse struct {
	ID              string              `json:"id"`
	UserID          string              `json:"user_id"`
	AssetID         *string             `json:"asset_id"`
	Content         string              `json:"content"`
	SentimentLabel  string              `json:"sentiment_label"`
	IsVisible       bool                `json:"is_visible"`
	LikesCount      int                 `json:"likes_count"`
	CreatedAt       time.Time           `json:"created_at"`
	Photos          []PostPhotoResponse `json:"photos"`
	AuthorName      string              `json:"author_name"`
	AuthorAvatarURL string              `json:"author_avatar_url"`
	CommentsCount   int                 `json:"comments_count"`
}

func FromEntity(p entities.Post, photos []entities.PostPhoto) PostResponse {
	photoResponses := make([]PostPhotoResponse, 0, len(photos))
	for _, ph := range photos {
		photoResponses = append(photoResponses, PostPhotoResponse{ID: ph.ID, URL: ph.URL, Order: ph.Order})
	}

	return PostResponse{
		ID:              p.ID,
		UserID:          p.UserID,
		AssetID:         p.AssetID,
		Content:         p.Content,
		SentimentLabel:  p.SentimentLabel,
		IsVisible:       p.IsVisible,
		LikesCount:      p.LikesCount,
		CreatedAt:       p.CreatedAt,
		Photos:          photoResponses,
		AuthorName:      p.AuthorName,
		AuthorAvatarURL: p.AuthorAvatarURL,
		CommentsCount:   p.CommentsCount,
	}
}

func FromEntities(list []entities.Post) []PostResponse {
	out := make([]PostResponse, 0, len(list))
	for _, p := range list {
		out = append(out, FromEntity(p, nil))
	}
	return out
}
