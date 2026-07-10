package application

import (
	"context"

	"vault/src/features/notifications/domain/dto/request"
	"vault/src/features/notifications/domain/dto/response"
	"vault/src/features/notifications/domain/entities"
	"vault/src/features/notifications/domain/repositories"
)

type CreateNotificationUseCase struct {
	repo repositories.NotificationRepository
}

func NewCreateNotificationUseCase(repo repositories.NotificationRepository) *CreateNotificationUseCase {
	return &CreateNotificationUseCase{repo: repo}
}

func (uc *CreateNotificationUseCase) Execute(ctx context.Context, userID string, req request.CreateNotificationRequest) (response.NotificationResponse, error) {
	if err := req.Validate(); err != nil {
		return response.NotificationResponse{}, err
	}

	created, err := uc.repo.Create(ctx, entities.Notification{
		UserID:  userID,
		Type:    req.Type,
		Subtype: req.Subtype,
		Title:   req.Title,
		Body:    req.Body,
		Data:    req.Data,
	})
	if err != nil {
		return response.NotificationResponse{}, err
	}

	return response.FromEntity(created), nil
}
