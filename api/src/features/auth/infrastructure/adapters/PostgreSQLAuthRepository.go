package adapters

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"vault/src/features/auth/domain/entities"
	"vault/src/features/auth/domain/repositories"
)

type PostgreSQLAuthRepository struct {
	pool *pgxpool.Pool
}

func NewPostgreSQLAuthRepository(pool *pgxpool.Pool) *PostgreSQLAuthRepository {
	return &PostgreSQLAuthRepository{pool: pool}
}

func (r *PostgreSQLAuthRepository) FindCredentialsByEmail(ctx context.Context, email string) (entities.Credentials, error) {
	const query = `
		SELECT id, name, email, password, COALESCE(avatar_url, ''), role
		FROM users
		WHERE email = $1
	`

	var c entities.Credentials
	err := r.pool.QueryRow(ctx, query, email).Scan(
		&c.UserID, &c.Name, &c.Email, &c.PasswordHash, &c.AvatarURL, &c.Role,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return entities.Credentials{}, repositories.ErrCredentialsNotFound
	}
	if err != nil {
		return entities.Credentials{}, fmt.Errorf("no se pudieron obtener las credenciales: %w", err)
	}

	return c, nil
}
