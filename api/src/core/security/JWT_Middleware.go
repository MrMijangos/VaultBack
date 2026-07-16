package security

import (
	"context"
	"net/http"
	"strings"
)

type contextKey string

const claimsContextKey contextKey = "claims"

// tokenFromRequest busca el JWT primero en el header Authorization: Bearer
// (usado por la app móvil, que no maneja cookies) y si no está, cae a la
// cookie httpOnly (usada por un eventual cliente web).
func tokenFromRequest(r *http.Request) string {
	if header := r.Header.Get("Authorization"); header != "" {
		if after, ok := strings.CutPrefix(header, "Bearer "); ok {
			return after
		}
	}
	if cookie, err := r.Cookie(AuthCookieName); err == nil {
		return cookie.Value
	}
	return ""
}

func RequireAuth(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := tokenFromRequest(r)
			if token == "" {
				http.Error(w, `{"error":"no autenticado"}`, http.StatusUnauthorized)
				return
			}

			claims, err := ParseToken(token, secret)
			if err != nil {
				http.Error(w, `{"error":"sesion invalida o expirada"}`, http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), claimsContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func ClaimsFromContext(ctx context.Context) (*Claims, bool) {
	claims, ok := ctx.Value(claimsContextKey).(*Claims)
	return claims, ok
}
