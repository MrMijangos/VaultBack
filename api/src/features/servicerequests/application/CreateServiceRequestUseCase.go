package application

import (
	"context"
	"fmt"
	"log"

	assetrepositories "vault/src/features/assets/domain/repositories"
	businessrepositories "vault/src/features/businesses/domain/repositories"
	notifentities "vault/src/features/notifications/domain/entities"
	"vault/src/features/servicerequests/domain/dto/request"
	"vault/src/features/servicerequests/domain/dto/response"
	"vault/src/features/servicerequests/domain/entities"
	"vault/src/features/servicerequests/domain/repositories"
)

type CreateServiceRequestUseCase struct {
	repo         repositories.ServiceRequestRepository
	assetRepo    assetrepositories.AssetRepository
	businessRepo businessrepositories.BusinessRepository
	notifier     repositories.NotificationPublisher
}

func NewCreateServiceRequestUseCase(
	repo repositories.ServiceRequestRepository,
	assetRepo assetrepositories.AssetRepository,
	businessRepo businessrepositories.BusinessRepository,
	notifier repositories.NotificationPublisher,
) *CreateServiceRequestUseCase {
	return &CreateServiceRequestUseCase{repo: repo, assetRepo: assetRepo, businessRepo: businessRepo, notifier: notifier}
}

func (uc *CreateServiceRequestUseCase) Execute(ctx context.Context, ownerID string, req request.CreateServiceRequestRequest) (response.ServiceRequestResponse, error) {
	if err := req.Validate(); err != nil {
		return response.ServiceRequestResponse{}, err
	}

	asset, err := uc.assetRepo.FindByID(ctx, req.AssetID)
	if err != nil {
		return response.ServiceRequestResponse{}, err
	}
	if asset.UserID != ownerID {
		return response.ServiceRequestResponse{}, assetrepositories.ErrAssetNotFound
	}

	business, err := uc.businessRepo.FindByID(ctx, req.BusinessID)
	if err != nil {
		return response.ServiceRequestResponse{}, err
	}

	sr := &entities.ServiceRequest{
		AssetID:    req.AssetID,
		OwnerID:    ownerID,
		BusinessID: req.BusinessID,
		Type:       req.Type,
		Status:     entities.ServiceRequestStatusPendienteAceptacion,
	}
	if err := uc.repo.Create(ctx, sr); err != nil {
		return response.ServiceRequestResponse{}, err
	}

	if _, err := uc.notifier.Create(ctx, notifentities.Notification{
		UserID:  business.UserID,
		Type:    req.Type,
		Subtype: "solicitud_recibida",
		Title:   "Nueva solicitud de " + typeLabel(req.Type),
		Body:    fmt.Sprintf("%s quiere que le hagas %s %s a %q.", asset.SellerName, articleFor(req.Type), typeLabel(req.Type), asset.Name),
	}); err != nil {
		log.Printf("no se pudo crear la notificacion de solicitud recibida: %v", err)
	}

	created, err := uc.repo.FindByID(ctx, sr.ID)
	if err != nil {
		return response.ServiceRequestResponse{}, err
	}
	return response.FromEntity(*created), nil
}
