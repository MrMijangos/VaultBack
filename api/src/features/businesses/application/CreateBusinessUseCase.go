package application

import (
	"context"

	"vault/src/features/businesses/domain/dto/request"
	"vault/src/features/businesses/domain/dto/response"
	"vault/src/features/businesses/domain/entities"
	"vault/src/features/businesses/domain/repositories"
)

type CreateBusinessUseCase struct {
	repo repositories.BusinessRepository
}

func NewCreateBusinessUseCase(repo repositories.BusinessRepository) *CreateBusinessUseCase {
	return &CreateBusinessUseCase{repo: repo}
}

func (uc *CreateBusinessUseCase) Execute(ctx context.Context, userID string, req request.CreateBusinessRequest) (response.BusinessResponse, error) {
	if err := req.Validate(); err != nil {
		return response.BusinessResponse{}, err
	}

	exists, err := uc.repo.ExistsByUserID(ctx, userID)
	if err != nil {
		return response.BusinessResponse{}, err
	}
	if exists {
		return response.BusinessResponse{}, repositories.ErrBusinessAlreadyExists
	}

	created, err := uc.repo.Create(ctx, entities.Business{
		UserID:      userID,
		Name:        req.Name,
		Type:        req.Type,
		Description: req.Description,
		Location:    req.Location,
		Specialties: req.Specialties,
	})
	if err != nil {
		return response.BusinessResponse{}, err
	}

	return response.FromEntity(created), nil
}
