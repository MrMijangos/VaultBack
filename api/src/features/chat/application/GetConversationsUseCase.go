package application

import (
	"context"

	"vault/src/features/chat/domain/dto/response"
	"vault/src/features/chat/domain/repositories"
)

type GetConversationsUseCase struct {
	repo repositories.ChatMessageRepository
}

func NewGetConversationsUseCase(repo repositories.ChatMessageRepository) *GetConversationsUseCase {
	return &GetConversationsUseCase{repo: repo}
}

func (uc *GetConversationsUseCase) Execute(ctx context.Context, userID string) ([]response.ConversationSummaryResponse, error) {
	list, err := uc.repo.FindConversationsForUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	return response.FromConversationSummaries(list), nil
}
