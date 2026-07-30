package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// loginLimiter le da a cada IP hasta 5 intentos de login, recargando 1 cada
// 3 minutos -- suficiente para frenar fuerza bruta automatizada sin
// molestar a alguien que de verdad se equivocó de contraseña un par de
// veces. Es en memoria: se reinicia si el servicio reinicia y no se
// comparte entre instancias, pero alcanza mientras corra una sola (ver
// Dockerfile/despliegue actual) sin sumar una dependencia como Redis.
type loginLimiter struct {
	mu       sync.Mutex
	visitors map[string]*rate.Limiter
}

func newLoginLimiter() *loginLimiter {
	return &loginLimiter{visitors: make(map[string]*rate.Limiter)}
}

func (l *loginLimiter) get(ip string) *rate.Limiter {
	l.mu.Lock()
	defer l.mu.Unlock()

	limiter, ok := l.visitors[ip]
	if !ok {
		limiter = rate.NewLimiter(rate.Every(3*time.Minute), 5)
		l.visitors[ip] = limiter
	}
	return limiter
}

// clientIP prioriza X-Real-IP porque en producción todo pasa por el nginx
// de gateway/ (ver gateway/nginx.conf), que sobreescribe RemoteAddr con su
// propia IP interna -- sin esto, todas las peticiones compartirían un único
// límite.
func clientIP(r *http.Request) string {
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	if host, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		return host
	}
	return r.RemoteAddr
}

func RateLimitLogin(next http.Handler) http.Handler {
	limiter := newLoginLimiter()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.get(clientIP(r)).Allow() {
			http.Error(w, `{"error":"demasiados intentos, intenta de nuevo en unos minutos"}`, http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
