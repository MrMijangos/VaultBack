package application

import (
	"context"

	"vault/src/features/chat/domain/dto/request"
	"vault/src/features/chat/domain/dto/response"
	"vault/src/features/chat/domain/repositories"
)

type UpdateChatMessageStatusUseCase struct {
	repo repositories.ChatMessageRepository
}

func NewUpdateChatMessageStatusUseCase(repo repositories.ChatMessageRepository) *UpdateChatMessageStatusUseCase {
	return &UpdateChatMessageStatusUseCase{repo: repo}
}

func (uc *UpdateChatMessageStatusUseCase) Execute(ctx context.Context, id string, recipientID string, req request.UpdateChatMessageStatusRequest) (response.ChatMessageResponse, error) {
	if err := req.Validate(); err != nil {
		return response.ChatMessageResponse{}, err
	}

	updated, err := uc.repo.UpdateStatus(ctx, id, recipientID, req.Status)
	if err != nil {
		return response.ChatMessageResponse{}, err
	}

	return response.FromEntity(updated), nil
}
