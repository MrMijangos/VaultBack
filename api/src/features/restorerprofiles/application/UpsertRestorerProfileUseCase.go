package application

import (
	"context"

	"vault/src/features/restorerprofiles/domain/dto/request"
	"vault/src/features/restorerprofiles/domain/dto/response"
	"vault/src/features/restorerprofiles/domain/entities"
	"vault/src/features/restorerprofiles/domain/repositories"
)

type UpsertRestorerProfileUseCase struct {
	repo repositories.RestorerProfileRepository
}

func NewUpsertRestorerProfileUseCase(repo repositories.RestorerProfileRepository) *UpsertRestorerProfileUseCase {
	return &UpsertRestorerProfileUseCase{repo: repo}
}

func (uc *UpsertRestorerProfileUseCase) Execute(ctx context.Context, userID string, req request.UpsertRestorerProfileRequest) (response.RestorerProfileResponse, error) {
	if err := req.Validate(); err != nil {
		return response.RestorerProfileResponse{}, err
	}

	services := make([]entities.RestorerService, 0, len(req.Services))
	for _, s := range req.Services {
		services = append(services, entities.RestorerService{
			Title:       s.Title,
			Description: s.Description,
			Price:       s.Price,
		})
	}

	specialties := req.Specialties
	if specialties == nil {
		specialties = []string{}
	}

	updated, err := uc.repo.Upsert(ctx, userID, req.Bio, specialties, services)
	if err != nil {
		return response.RestorerProfileResponse{}, err
	}

	return response.FromEntity(updated), nil
}
