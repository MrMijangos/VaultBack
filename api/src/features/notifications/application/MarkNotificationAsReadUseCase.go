package application

import (
	"context"

	"vault/src/features/notifications/domain/dto/response"
	"vault/src/features/notifications/domain/repositories"
)

type MarkNotificationAsReadUseCase struct {
	repo repositories.NotificationRepository
}

func NewMarkNotificationAsReadUseCase(repo repositories.NotificationRepository) *MarkNotificationAsReadUseCase {
	return &MarkNotificationAsReadUseCase{repo: repo}
}

func (uc *MarkNotificationAsReadUseCase) Execute(ctx context.Context, id string, userID string) (response.NotificationResponse, error) {
	updated, err := uc.repo.MarkAsRead(ctx, id, userID)
	if err != nil {
		return response.NotificationResponse{}, err
	}
	return response.FromEntity(updated), nil
}
