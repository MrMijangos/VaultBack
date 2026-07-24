package application

import (
	"context"

	"vault/src/features/businesses/domain/dto/request"
	"vault/src/features/businesses/domain/dto/response"
	"vault/src/features/businesses/domain/entities"
	"vault/src/features/businesses/domain/repositories"
)

type UpdateBusinessUseCase struct {
	repo repositories.BusinessRepository
}

func NewUpdateBusinessUseCase(repo repositories.BusinessRepository) *UpdateBusinessUseCase {
	return &UpdateBusinessUseCase{repo: repo}
}

func (uc *UpdateBusinessUseCase) Execute(ctx context.Context, id string, userID string, req request.UpdateBusinessRequest) (response.BusinessResponse, error) {
	if err := req.Validate(); err != nil {
		return response.BusinessResponse{}, err
	}

	updated, err := uc.repo.Update(ctx, id, userID, entities.Business{
		Name:        req.Name,
		Types:       req.Types,
		Description: req.Description,
		Location:    req.Location,
		Specialties: req.Specialties,
	})
	if err != nil {
		return response.BusinessResponse{}, err
	}

	return response.FromEntity(updated), nil
}
