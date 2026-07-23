package repositories

import (
	"context"
	"errors"

	"vault/src/features/chat/domain/entities"
)

var ErrChatMessageNotFound = errors.New("el mensaje no existe")

type ChatMessageRepository interface {
	Create(ctx context.Context, message entities.ChatMessage) (entities.ChatMessage, error)
	// FindConversation devuelve los mensajes entre userA y userB (en
	// cualquier direccion), ordenados por fecha de creacion.
	FindConversation(ctx context.Context, userA string, userB string) ([]entities.ChatMessage, error)
	// UpdateStatus solo puede ser invocado por el destinatario del mensaje.
	UpdateStatus(ctx context.Context, id string, recipientID string, status string) (entities.ChatMessage, error)
}
