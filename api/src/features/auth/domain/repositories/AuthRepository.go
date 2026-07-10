package repositories

import (
	"context"
	"errors"

	"vault/src/features/auth/domain/entities"
)

var ErrCredentialsNotFound = errors.New("no existe una cuenta con ese correo")

type AuthRepository interface {
	FindCredentialsByEmail(ctx context.Context, email string) (entities.Credentials, error)
}
