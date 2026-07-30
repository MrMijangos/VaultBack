package push

import (
	"context"
	"fmt"
	"log"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

// FCMSender es la implementación real, respaldada por la API HTTP v1 de FCM
// vía el SDK de administración de Firebase.
type FCMSender struct {
	client *messaging.Client
	tokens TokenProvider
}

// NewFCMSender inicializa el cliente de FCM a partir del JSON completo de la
// cuenta de servicio (ver FIREBASE_SERVICE_ACCOUNT_KEY en config.Config). Se
// llama una sola vez al arrancar el servidor y la instancia se reutiliza --
// *messaging.Client es seguro para uso concurrente.
func NewFCMSender(ctx context.Context, serviceAccountJSON string, tokens TokenProvider) (*FCMSender, error) {
	app, err := firebase.NewApp(ctx, nil, option.WithCredentialsJSON([]byte(serviceAccountJSON)))
	if err != nil {
		return nil, fmt.Errorf("no se pudo inicializar la app de Firebase: %w", err)
	}

	client, err := app.Messaging(ctx)
	if err != nil {
		return nil, fmt.Errorf("no se pudo inicializar el cliente de FCM: %w", err)
	}

	return &FCMSender{client: client, tokens: tokens}, nil
}

func (s *FCMSender) Notify(ctx context.Context, userID string, title string, body string, data map[string]string) error {
	tokens, err := s.tokens.FindTokensByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("no se pudieron obtener los tokens FCM de %s: %w", userID, err)
	}
	if len(tokens) == 0 {
		return nil
	}

	resp, err := s.client.SendEachForMulticast(ctx, &messaging.MulticastMessage{
		Tokens: tokens,
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Data: data,
	})
	if err != nil {
		return fmt.Errorf("no se pudo enviar el push a %s: %w", userID, err)
	}

	// Un token deja de servir cuando el usuario desinstala la app o el token
	// rota -- limpiarlo evita reintentar contra el mismo token en cada
	// notificación futura (ver sección 7 del doc de requisitos del back).
	for i, r := range resp.Responses {
		if r.Success {
			continue
		}
		if messaging.IsRegistrationTokenNotRegistered(r.Error) || messaging.IsUnregistered(r.Error) {
			if err := s.tokens.DeleteToken(ctx, tokens[i]); err != nil {
				log.Printf("[push] no se pudo limpiar el token FCM invalido: %v", err)
			}
		} else {
			log.Printf("[push] no se pudo enviar el push a un dispositivo de %s: %v", userID, r.Error)
		}
	}

	return nil
}
