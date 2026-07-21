package security

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

// Claims debe coincidir exactamente con security.Claims de api/ -- los
// tokens los emite api al hacer login, este servicio solo los valida.
// Mismo JWT_SECRET en ambos servicios (ver LoadConfig.go).
type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

var ErrInvalidToken = errors.New("token invalido o expirado")

func ParseToken(tokenString string, secret string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
