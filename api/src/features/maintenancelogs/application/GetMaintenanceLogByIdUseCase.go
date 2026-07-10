package application

import (
	"context"

	"vault/src/features/maintenancelogs/domain/dto/response"
	"vault/src/features/maintenancelogs/domain/repositories"
)

type GetMaintenanceLogByIdUseCase struct {
	repo repositories.MaintenanceLogRepository
}

func NewGetMaintenanceLogByIdUseCase(repo repositories.MaintenanceLogRepository) *GetMaintenanceLogByIdUseCase {
	return &GetMaintenanceLogByIdUseCase{repo: repo}
}

func (uc *GetMaintenanceLogByIdUseCase) Execute(ctx context.Context, id string) (response.MaintenanceLogResponse, error) {
	l, err := uc.repo.FindByID(ctx, id)
	if err != nil {
		return response.MaintenanceLogResponse{}, err
	}
	return response.FromEntity(l), nil
}
