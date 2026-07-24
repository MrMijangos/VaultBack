package repositories

import (
	"context"
	"errors"

	"vault/src/features/auth/domain/entities"
)

var ErrCredentialsNotFound = errors.New("no existe una cuenta con ese correo")

type AuthRepository interface {
	FindCredentialsByEmail(ctx context.Context, email string) (entities.Credentials, error)
	// FindOrCreateByEmail se usa para el login con Google: si ya existe una
	// cuenta con ese correo la reutiliza tal cual (name/avatarURL son los de
	// esa sesión de Google y no pisan lo que ya haya guardado); si no,
	// crea una cuenta nueva sin contraseña utilizable (el login por correo
	// para esa cuenta queda bloqueado hasta que registre una). isNew indica
	// si la cuenta se acaba de crear.
	FindOrCreateByEmail(ctx context.Context, email string, name string, avatarURL string) (entities.Credentials, bool, error)
}
