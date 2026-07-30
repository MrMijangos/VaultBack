package repositories

import "context"

type FCMTokenRepository interface {
	// Register hace upsert por token (UNIQUE en la tabla): si el token no
	// existe lo crea; si ya existe (mismo usuario o, tras una reinstalación,
	// otro usuario) lo reasigna a userID -- nunca duplica filas para el
	// mismo token.
	Register(ctx context.Context, userID string, token string, platform *string) error
	FindTokensByUserID(ctx context.Context, userID string) ([]string, error)
	// DeleteToken se llama tanto al hacer logout (endpoint DELETE) como desde
	// push.TokenProvider para limpiar tokens invalidos -- mismo nombre en
	// ambos puertos para que PostgreSQLFCMTokenRepository satisfaga los dos
	// por estructura (mismo patrón que AssetPhotoProvider en posts/).
	DeleteToken(ctx context.Context, token string) error
}
