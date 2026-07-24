package repositories

import (
	"context"
	"errors"

	"vault/src/features/users/domain/entities"
)

var ErrUserNotFound = errors.New("el usuario no existe")

type UserRepository interface {
	Create(ctx context.Context, user entities.User) (entities.User, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	FindAll(ctx context.Context) ([]entities.User, error)
	FindByID(ctx context.Context, id string) (entities.User, error)
	Update(ctx context.Context, id string, user entities.User) (entities.User, error)
	UpdateImage(ctx context.Context, id string, imageURL string) (entities.User, error)
	Delete(ctx context.Context, id string) error
	SetPublicKey(ctx context.Context, id string, publicKey string) error
	// GetPublicKey usa una consulta angosta (no reutiliza FindByID) para no
	// traer el hash de la contraseña en una ruta publica.
	GetPublicKey(ctx context.Context, id string) (*string, error)
	// AddRoles agrega roles al historico acumulado sin duplicar (dedup vía
	// SQL) y sin tocar los que ya tenia la cuenta.
	AddRoles(ctx context.Context, id string, roles []string) (entities.User, error)
}
