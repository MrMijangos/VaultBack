package application

import (
	"context"
	"fmt"
	"log"
	"time"

	businessrepositories "vault/src/features/businesses/domain/repositories"
	maintenancelogentities "vault/src/features/maintenancelogs/domain/entities"
	maintenancelogrepositories "vault/src/features/maintenancelogs/domain/repositories"
	notifentities "vault/src/features/notifications/domain/entities"
	"vault/src/features/servicerequests/domain/dto/response"
	"vault/src/features/servicerequests/domain/entities"
	"vault/src/features/servicerequests/domain/repositories"
)

// ConfirmServiceRequestUseCase lo llama el dueño del activo cuando ya
// recibió de vuelta su artículo -- cierra el flujo y, a diferencia de
// Accept/Start/Finish (que registran progreso), además certifica el
// trabajo creando la entrada correspondiente en maintenance_logs, así el
// contador de servicios/restauraciones de la card del activo (ver
// PostgreSQLAssetRepository.selectAssetsQuery) refleja trabajo real.
type ConfirmServiceRequestUseCase struct {
	repo            repositories.ServiceRequestRepository
	businessRepo    businessrepositories.BusinessRepository
	maintenanceLogs maintenancelogrepositories.MaintenanceLogRepository
	notifier        repositories.NotificationPublisher
}

func NewConfirmServiceRequestUseCase(
	repo repositories.ServiceRequestRepository,
	businessRepo businessrepositories.BusinessRepository,
	maintenanceLogs maintenancelogrepositories.MaintenanceLogRepository,
	notifier repositories.NotificationPublisher,
) *ConfirmServiceRequestUseCase {
	return &ConfirmServiceRequestUseCase{repo: repo, businessRepo: businessRepo, maintenanceLogs: maintenanceLogs, notifier: notifier}
}

func (uc *ConfirmServiceRequestUseCase) Execute(ctx context.Context, ownerID string, requestID string) (response.ServiceRequestResponse, error) {
	sr, err := uc.repo.FindByID(ctx, requestID)
	if err != nil {
		return response.ServiceRequestResponse{}, err
	}
	if sr == nil {
		return response.ServiceRequestResponse{}, repositories.ErrServiceRequestNotFound
	}
	if sr.OwnerID != ownerID {
		return response.ServiceRequestResponse{}, ErrNotOwner
	}
	if sr.Status != entities.ServiceRequestStatusTerminado {
		return response.ServiceRequestResponse{}, ErrInvalidTransition
	}

	business, err := uc.businessRepo.FindByID(ctx, sr.BusinessID)
	if err != nil {
		return response.ServiceRequestResponse{}, err
	}

	now := time.Now().UTC()
	sr.Status = entities.ServiceRequestStatusConfirmado
	sr.ConfirmedAt = &now
	if err := uc.repo.Update(ctx, sr); err != nil {
		return response.ServiceRequestResponse{}, err
	}

	maintenanceType := "mantenimiento"
	if sr.Type == entities.ServiceRequestTypeReparacion {
		maintenanceType = "restauracion"
	}
	providerID := business.UserID
	if _, err := uc.maintenanceLogs.Create(ctx, maintenancelogentities.MaintenanceLog{
		AssetID:     sr.AssetID,
		ProviderID:  &providerID,
		Type:        maintenanceType,
		PerformedAt: &now,
		Notes:       fmt.Sprintf("Generado automáticamente al confirmar la solicitud de %s con %s.", typeLabel(sr.Type), sr.BusinessName),
	}); err != nil {
		log.Printf("no se pudo registrar el mantenimiento automatico de la solicitud %s: %v", sr.ID, err)
	}

	if _, err := uc.notifier.Create(ctx, notifentities.Notification{
		UserID:  business.UserID,
		Type:    sr.Type,
		Subtype: "articulo_confirmado",
		Title:   "El cliente confirmó la entrega",
		Body:    fmt.Sprintf("%s confirmó que ya recibió su artículo de vuelta.", sr.OwnerName),
	}); err != nil {
		log.Printf("no se pudo crear la notificacion de articulo confirmado: %v", err)
	}

	updated, err := uc.repo.FindByID(ctx, sr.ID)
	if err != nil {
		return response.ServiceRequestResponse{}, err
	}
	return response.FromEntity(*updated), nil
}
