package application

import (
	"context"

	"vault/src/features/chat/domain/repositories"
)

type DeleteChatMessageUseCase struct {
	repo repositories.ChatMessageRepository
}

func NewDeleteChatMessageUseCase(repo repositories.ChatMessageRepository) *DeleteChatMessageUseCase {
	return &DeleteChatMessageUseCase{repo: repo}
}

func (uc *DeleteChatMessageUseCase) Execute(ctx context.Context, id string, userID string) error {
	return uc.repo.DeleteMessage(ctx, id, userID)
}
