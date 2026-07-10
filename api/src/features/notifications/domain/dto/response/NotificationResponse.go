package response

import (
	"encoding/json"
	"time"

	"vault/src/features/notifications/domain/entities"
)

type NotificationResponse struct {
	ID        string          `json:"id"`
	UserID    string          `json:"user_id"`
	Type      string          `json:"type"`
	Subtype   string          `json:"subtype"`
	Title     string          `json:"title"`
	Body      string          `json:"body"`
	Data      json.RawMessage `json:"data"`
	Read      bool            `json:"read"`
	CreatedAt time.Time       `json:"created_at"`
}

func FromEntity(n entities.Notification) NotificationResponse {
	return NotificationResponse{
		ID:        n.ID,
		UserID:    n.UserID,
		Type:      n.Type,
		Subtype:   n.Subtype,
		Title:     n.Title,
		Body:      n.Body,
		Data:      n.Data,
		Read:      n.Read,
		CreatedAt: n.CreatedAt,
	}
}

func FromEntities(list []entities.Notification) []NotificationResponse {
	out := make([]NotificationResponse, 0, len(list))
	for _, n := range list {
		out = append(out, FromEntity(n))
	}
	return out
}
