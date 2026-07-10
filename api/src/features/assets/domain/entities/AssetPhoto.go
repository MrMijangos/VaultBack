package entities

import "time"

type AssetPhoto struct {
	ID        string
	AssetID   string
	URL       string
	IsCover   bool
	Order     int
	CreatedAt time.Time
}
