package application

import (
	"context"

	"vault/src/features/notifications/domain/repositories"
)

type MarkAllNotificationsAsReadUseCase struct {
	repo repositories.NotificationRepository
}

func NewMarkAllNotificationsAsReadUseCase(repo repositories.NotificationRepository) *MarkAllNotificationsAsReadUseCase {
	return &MarkAllNotificationsAsReadUseCase{repo: repo}
}

func (uc *MarkAllNotificationsAsReadUseCase) Execute(ctx context.Context, userID string) error {
	return uc.repo.MarkAllAsRead(ctx, userID)
}
