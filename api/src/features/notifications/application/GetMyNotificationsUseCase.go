package application

import (
	"context"

	"vault/src/features/notifications/domain/dto/response"
	"vault/src/features/notifications/domain/repositories"
)

type GetMyNotificationsUseCase struct {
	repo repositories.NotificationRepository
}

func NewGetMyNotificationsUseCase(repo repositories.NotificationRepository) *GetMyNotificationsUseCase {
	return &GetMyNotificationsUseCase{repo: repo}
}

func (uc *GetMyNotificationsUseCase) Execute(ctx context.Context, userID string) ([]response.NotificationResponse, error) {
	list, err := uc.repo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return response.FromEntities(list), nil
}
