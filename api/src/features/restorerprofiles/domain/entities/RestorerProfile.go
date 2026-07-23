package entities

import "time"

type RestorerService struct {
	ID          string
	UserID      string
	Title       string
	Description string
	Price       float64
	CreatedAt   time.Time
}

type RestorerProfile struct {
	UserID       string
	Bio          string
	Specialties  []string
	Services     []RestorerService
	Rating       *float64
	TotalReviews int
	Name         string
	AvatarURL    string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
