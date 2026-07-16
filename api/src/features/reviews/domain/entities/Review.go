package entities

import "time"

type Review struct {
	ID              string
	UserID          string
	ProviderID      string
	Content         string
	SentimentScore  *float64
	ToxicityScore   *float64
	IsVisible       bool
	LikesCount      int
	CreatedAt       time.Time
	AuthorName      string
	AuthorAvatarURL string
}
