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
	PublicKey    *string
	// Roles es el historico acumulado de roles que la cuenta ha adquirido
	// (nunca se quita nada, solo se agrega) -- a diferencia de Role, que es
	// el mas reciente/principal.
	Roles []string
}
