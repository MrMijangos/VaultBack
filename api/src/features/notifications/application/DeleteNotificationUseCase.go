package application

import (
	"context"

	"vault/src/features/notifications/domain/repositories"
)

type DeleteNotificationUseCase struct {
	repo repositories.NotificationRepository
}

func NewDeleteNotificationUseCase(repo repositories.NotificationRepository) *DeleteNotificationUseCase {
	return &DeleteNotificationUseCase{repo: repo}
}

func (uc *DeleteNotificationUseCase) Execute(ctx context.Context, id string, userID string) error {
	return uc.repo.Delete(ctx, id, userID)
}
