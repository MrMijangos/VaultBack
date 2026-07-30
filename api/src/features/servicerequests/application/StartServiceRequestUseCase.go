package application

import (
	"context"
	"fmt"
	"log"
	"time"

	businessrepositories "vault/src/features/businesses/domain/repositories"
	notifentities "vault/src/features/notifications/domain/entities"
	"vault/src/features/servicerequests/domain/dto/response"
	"vault/src/features/servicerequests/domain/entities"
	"vault/src/features/servicerequests/domain/repositories"
)

// StartServiceRequestUseCase lo llama el negocio cuando empieza a trabajar
// el artículo (ya lo tenía en espera).
type StartServiceRequestUseCase struct {
	repo         repositories.ServiceRequestRepository
	businessRepo businessrepositories.BusinessRepository
	notifier     repositories.NotificationPublisher
}

func NewStartServiceRequestUseCase(
	repo repositories.ServiceRequestRepository,
	businessRepo businessrepositories.BusinessRepository,
	notifier repositories.NotificationPublisher,
) *StartServiceRequestUseCase {
	return &StartServiceRequestUseCase{repo: repo, businessRepo: businessRepo, notifier: notifier}
}

func (uc *StartServiceRequestUseCase) Execute(ctx context.Context, businessOwnerID string, requestID string) (response.ServiceRequestResponse, error) {
	sr, err := loadForBusinessAction(ctx, uc.repo, uc.businessRepo, businessOwnerID, requestID, entities.ServiceRequestStatusEnEspera)
	if err != nil {
		return response.ServiceRequestResponse{}, err
	}

	now := time.Now().UTC()
	sr.Status = entities.ServiceRequestStatusEnServicio
	sr.StartedAt = &now
	if err := uc.repo.Update(ctx, sr); err != nil {
		return response.ServiceRequestResponse{}, err
	}

	if _, err := uc.notifier.Create(ctx, notifentities.Notification{
		UserID:  sr.OwnerID,
		Type:    sr.Type,
		Subtype: subtypeFor("en_proceso", sr.Type),
		Title:   fmt.Sprintf("Tu artículo está en proceso de %s", typeLabel(sr.Type)),
		Body:    fmt.Sprintf("%s ya empezó a trabajar tu artículo.", sr.BusinessName),
	}); err != nil {
		log.Printf("no se pudo crear la notificacion de solicitud en proceso: %v", err)
	}

	updated, err := uc.repo.FindByID(ctx, sr.ID)
	if err != nil {
		return response.ServiceRequestResponse{}, err
	}
	return response.FromEntity(*updated), nil
}
