package push

import (
	"context"
	"log"
)

// Sender envía una notificación push a todos los dispositivos registrados de
// un usuario. Nunca debe tumbar el flujo que la dispara (crear un post,
// mandar un mensaje) -- si el usuario no tiene tokens, o si Firebase no está
// configurado, Notify simplemente no hace nada.
type Sender interface {
	Notify(ctx context.Context, userID string, title string, body string, data map[string]string) error
}

// TokenProvider es el puerto hacia fcmtokens/ -- mismo patrón que
// AssetPhotoProvider en posts/ o NotificationPublisher en chat/:
// PostgreSQLFCMTokenRepository lo satisface por estructura sin que este
// paquete dependa de su application/infraestructura.
type TokenProvider interface {
	FindTokensByUserID(ctx context.Context, userID string) ([]string, error)
	DeleteToken(ctx context.Context, token string) error
}

// NoopSender se usa cuando FIREBASE_SERVICE_ACCOUNT_KEY no está configurado
// o Firebase no pudo inicializarse: la API sigue funcionando con normalidad
// (las notificaciones dentro de la app se siguen creando), solo no se manda
// el push. Mismo criterio que eventbus.NoopPublisher.
type NoopSender struct{}

func (NoopSender) Notify(_ context.Context, userID string, title string, _ string, _ map[string]string) error {
	log.Printf("[push] Firebase no configurado, push %q para %s no enviado", title, userID)
	return nil
}
