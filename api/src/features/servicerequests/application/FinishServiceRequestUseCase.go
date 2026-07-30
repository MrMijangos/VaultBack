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

// FinishServiceRequestUseCase lo llama el negocio cuando termina el
// trabajo -- el dueño todavía tiene que confirmar que lo recibió de vuelta
// (ver ConfirmServiceRequestUseCase) para cerrar el flujo.
type FinishServiceRequestUseCase struct {
	repo         repositories.ServiceRequestRepository
	businessRepo businessrepositories.BusinessRepository
	notifier     repositories.NotificationPublisher
}

func NewFinishServiceRequestUseCase(
	repo repositories.ServiceRequestRepository,
	businessRepo businessrepositories.BusinessRepository,
	notifier repositories.NotificationPublisher,
) *FinishServiceRequestUseCase {
	return &FinishServiceRequestUseCase{repo: repo, businessRepo: businessRepo, notifier: notifier}
}

func (uc *FinishServiceRequestUseCase) Execute(ctx context.Context, businessOwnerID string, requestID string) (response.ServiceRequestResponse, error) {
	sr, err := loadForBusinessAction(ctx, uc.repo, uc.businessRepo, businessOwnerID, requestID, entities.ServiceRequestStatusEnServicio)
	if err != nil {
		return response.ServiceRequestResponse{}, err
	}

	now := time.Now().UTC()
	sr.Status = entities.ServiceRequestStatusTerminado
	sr.FinishedAt = &now
	if err := uc.repo.Update(ctx, sr); err != nil {
		return response.ServiceRequestResponse{}, err
	}

	if _, err := uc.notifier.Create(ctx, notifentities.Notification{
		UserID:  sr.OwnerID,
		Type:    sr.Type,
		Subtype: subtypeFor("salio", sr.Type),
		Title:   "Tu artículo está listo",
		Body:    fmt.Sprintf("%s ya terminó %s %s de tu artículo -- confírmalo cuando lo recibas de vuelta.", sr.BusinessName, definiteArticleFor(sr.Type), typeLabel(sr.Type)),
	}); err != nil {
		log.Printf("no se pudo crear la notificacion de solicitud terminada: %v", err)
	}

	updated, err := uc.repo.FindByID(ctx, sr.ID)
	if err != nil {
		return response.ServiceRequestResponse{}, err
	}
	return response.FromEntity(*updated), nil
}
