package security

import (
	"context"
	"net/http"
)

type contextKey string

const claimsContextKey contextKey = "claims"

func RequireAuth(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(AuthCookieName)
			if err != nil {
				http.Error(w, `{"error":"no autenticado"}`, http.StatusUnauthorized)
				return
			}

			claims, err := ParseToken(cookie.Value, secret)
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
