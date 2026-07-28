package application

import (
	"context"

	"vault/src/features/chat/domain/repositories"
)

type DeleteConversationUseCase struct {
	repo repositories.ChatMessageRepository
}

func NewDeleteConversationUseCase(repo repositories.ChatMessageRepository) *DeleteConversationUseCase {
	return &DeleteConversationUseCase{repo: repo}
}

func (uc *DeleteConversationUseCase) Execute(ctx context.Context, userID string, otherID string) error {
	return uc.repo.DeleteConversation(ctx, userID, otherID)
}
