package entities

import "time"

type Business struct {
	ID          string
	UserID      string
	Name        string
	Type        string
	Description string
	Location    string
	IsVerified  bool
	CreatedAt   time.Time
}
