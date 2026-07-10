package entities

import "time"

type PostPhoto struct {
	ID        string
	PostID    string
	URL       string
	Order     int
	CreatedAt time.Time
}
