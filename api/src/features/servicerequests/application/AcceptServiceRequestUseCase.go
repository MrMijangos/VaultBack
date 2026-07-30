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

// AcceptServiceRequestUseCase lo llama el negocio para confirmar que ya
// tiene el artículo en mano y puede empezar a trabajarlo.
type AcceptServiceRequestUseCase struct {
	repo         repositories.ServiceRequestRepository
	businessRepo businessrepositories.BusinessRepository
	notifier     repositories.NotificationPublisher
}

func NewAcceptServiceRequestUseCase(
	repo repositories.ServiceRequestRepository,
	businessRepo businessrepositories.BusinessRepository,
	notifier repositories.NotificationPublisher,
) *AcceptServiceRequestUseCase {
	return &AcceptServiceRequestUseCase{repo: repo, businessRepo: businessRepo, notifier: notifier}
}

func (uc *AcceptServiceRequestUseCase) Execute(ctx context.Context, businessOwnerID string, requestID string) (response.ServiceRequestResponse, error) {
	sr, err := loadForBusinessAction(ctx, uc.repo, uc.businessRepo, businessOwnerID, requestID, entities.ServiceRequestStatusPendienteAceptacion)
	if err != nil {
		return response.ServiceRequestResponse{}, err
	}

	now := time.Now().UTC()
	sr.Status = entities.ServiceRequestStatusEnEspera
	sr.AcceptedAt = &now
	if err := uc.repo.Update(ctx, sr); err != nil {
		return response.ServiceRequestResponse{}, err
	}

	if _, err := uc.notifier.Create(ctx, notifentities.Notification{
		UserID:  sr.OwnerID,
		Type:    sr.Type,
		Subtype: subtypeFor("entro", sr.Type),
		Title:   "Tu artículo está en espera",
		Body:    fmt.Sprintf("%s ya recibió tu artículo y está en espera para %s %s.", sr.BusinessName, definiteArticleFor(sr.Type), typeLabel(sr.Type)),
	}); err != nil {
		log.Printf("no se pudo crear la notificacion de solicitud aceptada: %v", err)
	}

	updated, err := uc.repo.FindByID(ctx, sr.ID)
	if err != nil {
		return response.ServiceRequestResponse{}, err
	}
	return response.FromEntity(*updated), nil
}
