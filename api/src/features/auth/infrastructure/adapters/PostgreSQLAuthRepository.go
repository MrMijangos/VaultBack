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
		SELECT id, name, email, password, COALESCE(avatar_url, ''), role, COALESCE(roles, '{}')
		FROM users
		WHERE email = $1
	`

	var c entities.Credentials
	err := r.pool.QueryRow(ctx, query, email).Scan(
		&c.UserID, &c.Name, &c.Email, &c.PasswordHash, &c.AvatarURL, &c.Role, &c.Roles,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return entities.Credentials{}, repositories.ErrCredentialsNotFound
	}
	if err != nil {
		return entities.Credentials{}, fmt.Errorf("no se pudieron obtener las credenciales: %w", err)
	}

	return c, nil
}

func (r *PostgreSQLAuthRepository) FindOrCreateByEmail(ctx context.Context, email string, name string, avatarURL string) (entities.Credentials, bool, error) {
	existing, err := r.FindCredentialsByEmail(ctx, email)
	if err == nil {
		return existing, false, nil
	}
	if !errors.Is(err, repositories.ErrCredentialsNotFound) {
		return entities.Credentials{}, false, err
	}

	// password = '' deja la cuenta sin login por correo posible (ver
	// security.VerifyPassword: un hash vacío nunca pasa el formato
	// argon2id esperado), y roles arranca con el mismo valor que role para
	// que ya tenga un primer rol acumulado como cualquier registro normal.
	const insertQuery = `
		INSERT INTO users (name, email, password, avatar_url, role, roles)
		VALUES ($1, $2, '', $3, 'usuario', ARRAY['usuario'])
		RETURNING id
	`
	var userID string
	if err := r.pool.QueryRow(ctx, insertQuery, name, email, avatarURL).Scan(&userID); err != nil {
		return entities.Credentials{}, false, fmt.Errorf("no se pudo crear la cuenta de Google: %w", err)
	}

	created, err := r.FindCredentialsByEmail(ctx, email)
	if err != nil {
		return entities.Credentials{}, false, err
	}
	return created, true, nil
}
