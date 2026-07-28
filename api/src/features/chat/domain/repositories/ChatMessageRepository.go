package repositories

import (
	"context"
	"errors"

	"vault/src/features/chat/domain/entities"
)

var ErrChatMessageNotFound = errors.New("el mensaje no existe")

type ChatMessageRepository interface {
	Create(ctx context.Context, message entities.ChatMessage) (entities.ChatMessage, error)
	// FindConversation devuelve los mensajes entre meID y otherID (en
	// cualquier direccion), ordenados por fecha de creacion, excluyendo los
	// que meID ya borró de su lado (ver DeleteMessage/DeleteConversation).
	FindConversation(ctx context.Context, meID string, otherID string) ([]entities.ChatMessage, error)
	// FindConversationsForUser devuelve un resumen por cada persona con la
	// que userID tiene al menos un mensaje sin borrar de su lado, con el más
	// reciente primero.
	FindConversationsForUser(ctx context.Context, userID string) ([]entities.ConversationSummary, error)
	// UpdateStatus solo puede ser invocado por el destinatario del mensaje.
	UpdateStatus(ctx context.Context, id string, recipientID string, status string) (entities.ChatMessage, error)
	// DeleteMessage borra el mensaje id solo del lado de userID (sender o
	// recipient) -- la otra persona lo sigue viendo. Falla con
	// ErrChatMessageNotFound si el mensaje no existe o userID no participa.
	DeleteMessage(ctx context.Context, id string, userID string) error
	// DeleteConversation borra todos los mensajes entre userID y otherID,
	// solo del lado de userID.
	DeleteConversation(ctx context.Context, userID string, otherID string) error
}
