package adapters

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"vault/src/features/users/domain/entities"
	"vault/src/features/users/domain/repositories"
)

const selectUsersQuery = `
	SELECT id, name, email, password, COALESCE(avatar_url, ''), role, created_at, updated_at
	FROM users
`

type PostgreSQLUserRepository struct {
	pool *pgxpool.Pool
}

func NewPostgreSQLUserRepository(pool *pgxpool.Pool) *PostgreSQLUserRepository {
	return &PostgreSQLUserRepository{pool: pool}
}

func (r *PostgreSQLUserRepository) Create(ctx context.Context, user entities.User) (entities.User, error) {
	const query = `
		INSERT INTO users (name, email, password, avatar_url, role)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	err := r.pool.QueryRow(ctx, query,
		user.Name, user.Email, user.PasswordHash, user.AvatarURL, user.Role,
	).Scan(&user.ID)
	if err != nil {
		return entities.User{}, fmt.Errorf("no se pudo crear el usuario: %w", err)
	}

	return r.FindByID(ctx, user.ID)
}

func (r *PostgreSQLUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	const query = `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`

	var exists bool
	if err := r.pool.QueryRow(ctx, query, email).Scan(&exists); err != nil {
		return false, fmt.Errorf("no se pudo verificar el correo: %w", err)
	}
	return exists, nil
}

func (r *PostgreSQLUserRepository) FindAll(ctx context.Context) ([]entities.User, error) {
	rows, err := r.pool.Query(ctx, selectUsersQuery+" ORDER BY created_at")
	if err != nil {
		return nil, fmt.Errorf("no se pudieron listar los usuarios: %w", err)
	}
	defer rows.Close()

	var users []entities.User
	for rows.Next() {
		var u entities.User
		if err := rows.Scan(
			&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.AvatarURL, &u.Role, &u.CreatedAt, &u.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("no se pudo leer el usuario: %w", err)
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error al listar usuarios: %w", err)
	}

	return users, nil
}

func (r *PostgreSQLUserRepository) FindByID(ctx context.Context, id string) (entities.User, error) {
	var u entities.User
	err := r.pool.QueryRow(ctx, selectUsersQuery+" WHERE id = $1", id).Scan(
		&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.AvatarURL, &u.Role, &u.CreatedAt, &u.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return entities.User{}, repositories.ErrUserNotFound
	}
	if err != nil {
		return entities.User{}, fmt.Errorf("no se pudo obtener el usuario: %w", err)
	}

	return u, nil
}

func (r *PostgreSQLUserRepository) Update(ctx context.Context, id string, user entities.User) (entities.User, error) {
	const query = `
		UPDATE users
		SET name = $1, avatar_url = $2, role = $3, updated_at = now()
		WHERE id = $4
	`

	tag, err := r.pool.Exec(ctx, query, user.Name, user.AvatarURL, user.Role, id)
	if err != nil {
		return entities.User{}, fmt.Errorf("no se pudo actualizar el usuario: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return entities.User{}, repositories.ErrUserNotFound
	}

	return r.FindByID(ctx, id)
}

func (r *PostgreSQLUserRepository) UpdateImage(ctx context.Context, id string, imageURL string) (entities.User, error) {
	const query = `UPDATE users SET avatar_url = $1, updated_at = now() WHERE id = $2`

	tag, err := r.pool.Exec(ctx, query, imageURL, id)
	if err != nil {
		return entities.User{}, fmt.Errorf("no se pudo actualizar la imagen del usuario: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return entities.User{}, repositories.ErrUserNotFound
	}

	return r.FindByID(ctx, id)
}

func (r *PostgreSQLUserRepository) SetPublicKey(ctx context.Context, id string, publicKey string) error {
	tag, err := r.pool.Exec(ctx, `UPDATE users SET public_key = $1 WHERE id = $2`, publicKey, id)
	if err != nil {
		return fmt.Errorf("no se pudo guardar la llave publica: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return repositories.ErrUserNotFound
	}
	return nil
}

func (r *PostgreSQLUserRepository) GetPublicKey(ctx context.Context, id string) (*string, error) {
	var publicKey *string
	err := r.pool.QueryRow(ctx, `SELECT public_key FROM users WHERE id = $1`, id).Scan(&publicKey)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, repositories.ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("no se pudo obtener la llave publica: %w", err)
	}
	return publicKey, nil
}

func (r *PostgreSQLUserRepository) Delete(ctx context.Context, id string) error {
	const query = `DELETE FROM users WHERE id = $1`

	tag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("no se pudo eliminar el usuario: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return repositories.ErrUserNotFound
	}

	return nil
}
