package security

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthCookieName debe coincidir con api/src/core/security/Utils.go -- la
// app móvil manda el JWT por header Authorization, pero se soporta también
// la cookie httpOnly por si algún día hay un cliente web.
const AuthCookieName = "vault_token"

const ClaimsContextKey = "claims"

func tokenFromRequest(c *gin.Context) string {
	if header := c.GetHeader("Authorization"); header != "" {
		if after, ok := strings.CutPrefix(header, "Bearer "); ok {
			return after
		}
	}
	if cookie, err := c.Cookie(AuthCookieName); err == nil {
		return cookie
	}
	return ""
}

func RequireAuth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := tokenFromRequest(c)
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "no autenticado"})
			return
		}

		claims, err := ParseToken(token, secret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "sesion invalida o expirada"})
			return
		}

		c.Set(ClaimsContextKey, claims)
		c.Next()
	}
}

func ClaimsFromContext(c *gin.Context) (*Claims, bool) {
	value, ok := c.Get(ClaimsContextKey)
	if !ok {
		return nil, false
	}
	claims, ok := value.(*Claims)
	return claims, ok
}
