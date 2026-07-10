package repositories

import (
	"context"
	"errors"

	"vault/src/features/maintenancelogs/domain/entities"
)

var ErrMaintenanceLogNotFound = errors.New("el registro de servicio no existe")

type MaintenanceLogRepository interface {
	Create(ctx context.Context, log entities.MaintenanceLog) (entities.MaintenanceLog, error)
	FindByAssetID(ctx context.Context, assetID string) ([]entities.MaintenanceLog, error)
	FindByID(ctx context.Context, id string) (entities.MaintenanceLog, error)
	Update(ctx context.Context, id string, log entities.MaintenanceLog) (entities.MaintenanceLog, error)
	Delete(ctx context.Context, id string) error
}
