package application

import (
	"context"

	"vault/src/features/chat/domain/dto/request"
	"vault/src/features/chat/domain/dto/response"
	"vault/src/features/chat/domain/entities"
	"vault/src/features/chat/domain/repositories"
)

// SendChatMessageUseCase no llama a moderacion NLP -- el contenido es
// cifrado extremo a extremo, el servidor nunca ve texto plano y no puede
// (ni debe) analizarlo.
type SendChatMessageUseCase struct {
	repo repositories.ChatMessageRepository
}

func NewSendChatMessageUseCase(repo repositories.ChatMessageRepository) *SendChatMessageUseCase {
	return &SendChatMessageUseCase{repo: repo}
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

	return response.FromEntity(created), nil
}
