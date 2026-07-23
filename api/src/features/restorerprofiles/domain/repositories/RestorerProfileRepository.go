package repositories

import (
	"context"
	"errors"

	"vault/src/features/restorerprofiles/domain/entities"
)

var ErrProfileNotFound = errors.New("el perfil no existe")

type RestorerProfileRepository interface {
	// Upsert reemplaza por completo el perfil (bio + specialties) y su
	// lista de servicios, en una sola transaccion.
	Upsert(ctx context.Context, userID string, bio string, specialties []string, services []entities.RestorerService) (entities.RestorerProfile, error)
	FindByUserID(ctx context.Context, userID string) (entities.RestorerProfile, error)
	// ListWithServices solo devuelve perfiles que tienen al menos un servicio.
	ListWithServices(ctx context.Context) ([]entities.RestorerProfile, error)
}
