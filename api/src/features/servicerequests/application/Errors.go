package application

import "errors"

var (
	ErrNotOwner          = errors.New("no eres el dueño de este activo")
	ErrNotBusinessOwner  = errors.New("no eres el dueño de este negocio")
	ErrInvalidTransition = errors.New("la solicitud no está en un estado válido para esa acción")
)
