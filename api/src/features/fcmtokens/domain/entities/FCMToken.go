package entities

import "time"

type FCMToken struct {
	ID        string
	UserID    string
	Token     string
	Platform  *string
	CreatedAt time.Time
	UpdatedAt time.Time
}
