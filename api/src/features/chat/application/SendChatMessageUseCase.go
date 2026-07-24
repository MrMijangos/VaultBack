package application

import (
	"context"
	"log"

	"vault/src/features/chat/domain/dto/request"
	"vault/src/features/chat/domain/dto/response"
	"vault/src/features/chat/domain/entities"
	"vault/src/features/chat/domain/repositories"
	notifentities "vault/src/features/notifications/domain/entities"
)

// SendChatMessageUseCase no llama a moderacion NLP -- el contenido es
// cifrado extremo a extremo, el servidor nunca ve texto plano y no puede
// (ni debe) analizarlo.
type SendChatMessageUseCase struct {
	repo     repositories.ChatMessageRepository
	notifier repositories.NotificationPublisher
}

func NewSendChatMessageUseCase(repo repositories.ChatMessageRepository, notifier repositories.NotificationPublisher) *SendChatMessageUseCase {
	return &SendChatMessageUseCase{repo: repo, notifier: notifier}
}

func (uc *SendChatMessageUseCase) Execute(ctx context.Context, senderID string, req request.SendChatMessageRequest) (response.ChatMessageResponse, error) {
	if err := req.Validate(); err != nil {
		return response.ChatMessageResponse{}, err
	}

	created, err := uc.repo.Create(ctx, entities.ChatMessage{
		SenderID:              senderID,
		RecipientID:           req.RecipientID,
		CipherText:            req.CipherText,
		EncryptedAESKey:       req.EncryptedAESKey,
		EncryptedAESKeySender: req.EncryptedAESKeySender,
		IV:                    req.IV,
		Status:                "sent",
	})
	if err != nil {
		return response.ChatMessageResponse{}, err
	}

	// Si la notificación falla, el mensaje ya se mandó igual -- no vale la
	// pena tumbar el envío por esto, solo se registra.
	if _, err := uc.notifier.Create(ctx, notifentities.Notification{
		UserID:  req.RecipientID,
		Type:    "mensaje",
		Subtype: "mensaje_nuevo",
		Title:   "Nuevo mensaje",
		Body:    "Tienes un mensaje nuevo en el chat.",
	}); err != nil {
		log.Printf("no se pudo crear la notificacion de mensaje nuevo: %v", err)
	}

	return response.FromEntity(created), nil
}
