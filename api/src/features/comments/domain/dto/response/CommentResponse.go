package response

import (
	"time"

	"vault/src/features/comments/domain/entities"
)

type CommentResponse struct {
	ID        string    `json:"id"`
	PostID    string    `json:"post_id"`
	UserID    string    `json:"user_id"`
	Content   string    `json:"content"`
	IsVisible bool      `json:"is_visible"`
	CreatedAt time.Time `json:"created_at"`
}

func FromEntity(c entities.Comment) CommentResponse {
	return CommentResponse{
		ID:        c.ID,
		PostID:    c.PostID,
		UserID:    c.UserID,
		Content:   c.Content,
		IsVisible: c.IsVisible,
		CreatedAt: c.CreatedAt,
	}
}

func FromEntities(list []entities.Comment) []CommentResponse {
	out := make([]CommentResponse, 0, len(list))
	for _, c := range list {
		out = append(out, FromEntity(c))
	}
	return out
}
