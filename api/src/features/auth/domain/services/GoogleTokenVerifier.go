package services

import "context"

// GoogleProfile es lo que se obtiene de un idToken de Google ya verificado.
type GoogleProfile struct {
	Email   string
	Name    string
	Picture string
}

// GoogleTokenVerifier es el puerto que el caso de uso de login con Google
// necesita -- la implementación real (infrastructure/adapters) decide cómo
// se verifica el token (hoy: el endpoint tokeninfo de Google).
type GoogleTokenVerifier interface {
	Verify(ctx context.Context, idToken string) (GoogleProfile, error)
}
