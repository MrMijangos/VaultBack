package middleware

import "net/http"

func CORS(allowedOrigin string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
			// "*" + Allow-Credentials:true es una combinación inválida (el
			// navegador la rechaza para cualquier request con cookies) --
			// sin CORS_ORIGIN configurado en Railway, cfg.CORSOrigin cae a
			// "*" (ver LoadConfig.go) y este header nunca debe mandarse en
			// ese caso, o rompe en silencio el login por cookie desde un
			// cliente web/navegador.
			if allowedOrigin != "*" {
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
