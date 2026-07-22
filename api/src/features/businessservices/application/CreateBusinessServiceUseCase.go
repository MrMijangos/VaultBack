package application

import (
	"context"

	"vault/src/features/businessservices/domain/dto/request"
	"vault/src/features/businessservices/domain/dto/response"
	"vault/src/features/businessservices/domain/entities"
	"vault/src/features/businessservices/domain/repositories"
)

type CreateBusinessServiceUseCase struct {
	repo          repositories.BusinessServiceRepository
	ownerProvider repositories.BusinessOwnerProvider
}

func NewCreateBusinessServiceUseCase(repo repositories.BusinessServiceRepository, ownerProvider repositories.BusinessOwnerProvider) *CreateBusinessServiceUseCase {
	return &CreateBusinessServiceUseCase{repo: repo, ownerProvider: ownerProvider}
}

func (uc *CreateBusinessServiceUseCase) Execute(ctx context.Context, businessID string, userID string, req request.BusinessServiceRequest) (response.BusinessServiceResponse, error) {
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

	created, err := uc.repo.Create(ctx, entities.BusinessService{
		BusinessID:  businessID,
		Title:       req.Title,
		Description: req.Description,
		Price:       req.Price,
	})
	if err != nil {
		return response.BusinessServiceResponse{}, err
	}

	return response.FromEntity(created), nil
}
