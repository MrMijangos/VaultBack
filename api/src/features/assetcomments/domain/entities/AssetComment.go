package entities

import "time"

type AssetComment struct {
	ID              string
	AssetID         string
	UserID          string
	Content         string
	ToxicityScore   *float64
	IsVisible       bool
	CreatedAt       time.Time
	AuthorName      string
	AuthorAvatarURL string
}
