package application

import (
	"context"

	"vault/src/features/businessservices/domain/dto/request"
	"vault/src/features/businessservices/domain/dto/response"
	"vault/src/features/businessservices/domain/entities"
	"vault/src/features/businessservices/domain/repositories"
)

type UpdateBusinessServiceUseCase struct {
	repo          repositories.BusinessServiceRepository
	ownerProvider repositories.BusinessOwnerProvider
}

func NewUpdateBusinessServiceUseCase(repo repositories.BusinessServiceRepository, ownerProvider repositories.BusinessOwnerProvider) *UpdateBusinessServiceUseCase {
	return &UpdateBusinessServiceUseCase{repo: repo, ownerProvider: ownerProvider}
}

func (uc *UpdateBusinessServiceUseCase) Execute(ctx context.Context, businessID string, serviceID string, userID string, req request.BusinessServiceRequest) (response.BusinessServiceResponse, error) {
	if err := req.Validate(); err != nil {
		return response.BusinessServiceResponse{}, err
	}

	ownerID, err := uc.ownerProvider.GetOwnerUserID(ctx, businessID)
	if err != nil {
		return response.BusinessServiceResponse{}, err
	}
	if ownerID != userID {
		return response.BusinessServiceResponse{}, repositories.ErrNotOwner
	}

	updated, err := uc.repo.Update(ctx, serviceID, businessID, entities.BusinessService{
		Title:       req.Title,
		Description: req.Description,
		Price:       req.Price,
	})
	if err != nil {
		return response.BusinessServiceResponse{}, err
	}

	return response.FromEntity(updated), nil
}
