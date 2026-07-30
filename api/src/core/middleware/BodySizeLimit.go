package middleware

import "net/http"

// MaxRequestBodySize es el techo general para cualquier request -- sin este
// límite un JSON gigante se lee entero en memoria antes de fallar la
// validación, lo que permite tirar el servicio por RAM con pocas peticiones
// concurrentes. Las rutas de subida de fotos ya aplican su propio límite más
// chico encima de este (ver maxUploadSize en cada UploadXController).
const MaxRequestBodySize = 10 << 20 // 10MB

func LimitBodySize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, MaxRequestBodySize)
		next.ServeHTTP(w, r)
	})
}
