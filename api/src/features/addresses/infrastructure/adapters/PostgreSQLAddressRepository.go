package adapters

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"vault/src/features/addresses/domain/entities"
	"vault/src/features/addresses/domain/repositories"
)

const selectAddressesQuery = `
	SELECT id, user_id, label, recipient, phone, street, city, state, postal_code,
	       COALESCE(reference_notes, ''), is_default, created_at
	FROM addresses
`

type PostgreSQLAddressRepository struct {
	pool *pgxpool.Pool
}

func NewPostgreSQLAddressRepository(pool *pgxpool.Pool) *PostgreSQLAddressRepository {
	return &PostgreSQLAddressRepository{pool: pool}
}

func scanAddress(row pgx.Row) (entities.Address, error) {
	var a entities.Address
	err := row.Scan(
		&a.ID, &a.UserID, &a.Label, &a.Recipient, &a.Phone, &a.Street, &a.City, &a.State, &a.PostalCode,
		&a.ReferenceNotes, &a.IsDefault, &a.CreatedAt,
	)
	return a, err
}

func (r *PostgreSQLAddressRepository) Create(ctx context.Context, address entities.Address) (entities.Address, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return entities.Address{}, fmt.Errorf("no se pudo iniciar la transaccion: %w", err)
	}
	defer tx.Rollback(ctx)

	var count int
	if err := tx.QueryRow(ctx, `SELECT COUNT(*) FROM addresses WHERE user_id = $1`, address.UserID).Scan(&count); err != nil {
		return entities.Address{}, fmt.Errorf("no se pudo verificar las direcciones existentes: %w", err)
	}
	isDefault := count == 0

	const query = `
		INSERT INTO addresses (user_id, label, recipient, phone, street, city, state, postal_code, reference_notes, is_default)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, user_id, label, recipient, phone, street, city, state, postal_code, COALESCE(reference_notes, ''), is_default, created_at
	`
	row := tx.QueryRow(ctx, query,
		address.UserID, address.Label, address.Recipient, address.Phone, address.Street,
		address.City, address.State, address.PostalCode, address.ReferenceNotes, isDefault,
	)
	created, err := scanAddress(row)
	if err != nil {
		return entities.Address{}, fmt.Errorf("no se pudo crear la direccion: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return entities.Address{}, fmt.Errorf("no se pudo confirmar la transaccion: %w", err)
	}

	return created, nil
}

func (r *PostgreSQLAddressRepository) ListByUserID(ctx context.Context, userID string) ([]entities.Address, error) {
	rows, err := r.pool.Query(ctx, selectAddressesQuery+" WHERE user_id = $1 ORDER BY created_at", userID)
	if err != nil {
		return nil, fmt.Errorf("no se pudieron listar las direcciones: %w", err)
	}
	defer rows.Close()

	var list []entities.Address
	for rows.Next() {
		a, err := scanAddress(rows)
		if err != nil {
			return nil, fmt.Errorf("no se pudo leer la direccion: %w", err)
		}
		list = append(list, a)
	}
	return list, rows.Err()
}

func (r *PostgreSQLAddressRepository) Delete(ctx context.Context, id string, userID string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("no se pudo iniciar la transaccion: %w", err)
	}
	defer tx.Rollback(ctx)

	var wasDefault bool
	err = tx.QueryRow(ctx, `SELECT is_default FROM addresses WHERE id = $1 AND user_id = $2 FOR UPDATE`, id, userID).Scan(&wasDefault)
	if errors.Is(err, pgx.ErrNoRows) {
		return repositories.ErrAddressNotFound
	}
	if err != nil {
		return fmt.Errorf("no se pudo verificar la direccion: %w", err)
	}

	if _, err := tx.Exec(ctx, `DELETE FROM addresses WHERE id = $1 AND user_id = $2`, id, userID); err != nil {
		return fmt.Errorf("no se pudo eliminar la direccion: %w", err)
	}

	if wasDefault {
		if _, err := tx.Exec(ctx, `
			UPDATE addresses SET is_default = true WHERE id = (
				SELECT id FROM addresses WHERE user_id = $1 ORDER BY created_at LIMIT 1
			)
		`, userID); err != nil {
			return fmt.Errorf("no se pudo promover una nueva direccion predeterminada: %w", err)
		}
	}

	return tx.Commit(ctx)
}

func (r *PostgreSQLAddressRepository) SetDefault(ctx context.Context, id string, userID string) (entities.Address, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return entities.Address{}, fmt.Errorf("no se pudo iniciar la transaccion: %w", err)
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `UPDATE addresses SET is_default = false WHERE user_id = $1`, userID); err != nil {
		return entities.Address{}, fmt.Errorf("no se pudo actualizar las direcciones: %w", err)
	}

	tag, err := tx.Exec(ctx, `UPDATE addresses SET is_default = true WHERE id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		return entities.Address{}, fmt.Errorf("no se pudo marcar la direccion como predeterminada: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return entities.Address{}, repositories.ErrAddressNotFound
	}

	row := tx.QueryRow(ctx, selectAddressesQuery+" WHERE id = $1", id)
	updated, err := scanAddress(row)
	if err != nil {
		return entities.Address{}, fmt.Errorf("no se pudo leer la direccion actualizada: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return entities.Address{}, fmt.Errorf("no se pudo confirmar la transaccion: %w", err)
	}

	return updated, nil
}
