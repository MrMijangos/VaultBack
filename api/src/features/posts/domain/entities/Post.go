package entities

import "time"

type Post struct {
	ID             string
	UserID         string
	AssetID        *string
	Content        string
	SentimentScore *float64
	SentimentLabel string
	ToxicityScore  *float64
	IsVisible      bool
	LikesCount     int
	CreatedAt      time.Time
}
