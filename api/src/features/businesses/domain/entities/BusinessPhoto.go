package entities

import "time"

type BusinessPhoto struct {
	ID         string
	BusinessID string
	URL        string
	IsCover    bool
	Order      int
	CreatedAt  time.Time
}
