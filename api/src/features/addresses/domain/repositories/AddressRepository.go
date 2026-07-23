package repositories

import (
	"context"
	"errors"

	"vault/src/features/addresses/domain/entities"
)

var ErrAddressNotFound = errors.New("la direccion no existe")

type AddressRepository interface {
	// Create inserta la direccion; si es la primera del usuario, queda
	// is_default=true automaticamente.
	Create(ctx context.Context, address entities.Address) (entities.Address, error)
	ListByUserID(ctx context.Context, userID string) ([]entities.Address, error)
	// Delete borra la direccion del usuario; si era la predeterminada y
	// quedan otras, promueve la mas antigua restante.
	Delete(ctx context.Context, id string, userID string) error
	// SetDefault desmarca cualquier otra direccion default del usuario y
	// marca esta.
	SetDefault(ctx context.Context, id string, userID string) (entities.Address, error)
}
