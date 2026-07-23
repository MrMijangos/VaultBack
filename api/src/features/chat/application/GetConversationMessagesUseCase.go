package application

import (
	"context"

	"vault/src/features/chat/domain/dto/response"
	"vault/src/features/chat/domain/repositories"
)

type GetConversationMessagesUseCase struct {
	repo repositories.ChatMessageRepository
}

func NewGetConversationMessagesUseCase(repo repositories.ChatMessageRepository) *GetConversationMessagesUseCase {
	return &GetConversationMessagesUseCase{repo: repo}
}

func (uc *GetConversationMessagesUseCase) Execute(ctx context.Context, meID string, otherID string) ([]response.ChatMessageResponse, error) {
	list, err := uc.repo.FindConversation(ctx, meID, otherID)
	if err != nil {
		return nil, err
	}
	return response.FromEntities(list), nil
}
