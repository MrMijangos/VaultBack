package adapters

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"vault/src/features/auth/domain/services"
)

var ErrInvalidGoogleToken = errors.New("token de Google invalido")

// GoogleTokenVerifier valida un idToken de Google Sign-In contra el
// endpoint tokeninfo de Google -- Google hace la verificación criptográfica
// de la firma por nosotros, así que no hace falta manejar JWKS ni claves
// públicas en el backend. https://oauth2.googleapis.com/tokeninfo también
// devuelve 400 si el token expiró o es inválido.
type GoogleTokenVerifier struct {
	httpClient       *http.Client
	allowedAudiences map[string]bool
}

func NewGoogleTokenVerifier(allowedClientIDs string) *GoogleTokenVerifier {
	allowed := map[string]bool{}
	for _, id := range strings.Split(allowedClientIDs, ",") {
		id = strings.TrimSpace(id)
		if id != "" {
			allowed[id] = true
		}
	}
	return &GoogleTokenVerifier{
		httpClient:       &http.Client{Timeout: 5 * time.Second},
		allowedAudiences: allowed,
	}
}

func (v *GoogleTokenVerifier) Verify(ctx context.Context, idToken string) (services.GoogleProfile, error) {
	if len(v.allowedAudiences) == 0 {
		return services.GoogleProfile{}, errors.New("el inicio de sesion con Google no esta configurado (falta GOOGLE_OAUTH_CLIENT_IDS)")
	}

	endpoint := "https://oauth2.googleapis.com/tokeninfo?id_token=" + url.QueryEscape(idToken)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return services.GoogleProfile{}, err
	}

	resp, err := v.httpClient.Do(req)
	if err != nil {
		return services.GoogleProfile{}, fmt.Errorf("no se pudo verificar el token de Google: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return services.GoogleProfile{}, ErrInvalidGoogleToken
	}

	var payload struct {
		Audience      string `json:"aud"`
		Email         string `json:"email"`
		EmailVerified string `json:"email_verified"`
		Name          string `json:"name"`
		Picture       string `json:"picture"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return services.GoogleProfile{}, fmt.Errorf("respuesta invalida de Google: %w", err)
	}

	if !v.allowedAudiences[payload.Audience] {
		return services.GoogleProfile{}, ErrInvalidGoogleToken
	}
	// tokeninfo devuelve email_verified como string "true"/"false", no bool.
	if payload.Email == "" || payload.EmailVerified != "true" {
		return services.GoogleProfile{}, ErrInvalidGoogleToken
	}

	return services.GoogleProfile{
		Email:   payload.Email,
		Name:    payload.Name,
		Picture: payload.Picture,
	}, nil
}
