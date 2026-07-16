package entities

import "time"

type Comment struct {
	ID              string
	PostID          string
	UserID          string
	Content         string
	ToxicityScore   *float64
	IsVisible       bool
	CreatedAt       time.Time
	AuthorName      string
	AuthorAvatarURL string
}
