package entities

import "time"

type User struct {
	ID           string
	Name         string
	Email        string
	PasswordHash string
	AvatarURL    string
	Role         string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
