package adapters

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"vault/src/features/businesses/domain/entities"
	"vault/src/features/businesses/domain/repositories"
)

const selectBusinessesQuery = `
	SELECT id, user_id, name, type, COALESCE(description, ''), COALESCE(location, ''), is_verified, created_at, COALESCE(specialties, '{}')
	FROM businesses
`

type PostgreSQLBusinessRepository struct {
	pool *pgxpool.Pool
}

func NewPostgreSQLBusinessRepository(pool *pgxpool.Pool) *PostgreSQLBusinessRepository {
	return &PostgreSQLBusinessRepository{pool: pool}
}

func scanBusiness(row pgx.Row) (entities.Business, error) {
	var b entities.Business
	err := row.Scan(&b.ID, &b.UserID, &b.Name, &b.Type, &b.Description, &b.Location, &b.IsVerified, &b.CreatedAt, &b.Specialties)
	return b, err
}

func (r *PostgreSQLBusinessRepository) Create(ctx context.Context, business entities.Business) (entities.Business, error) {
	const query = `
		INSERT INTO businesses (user_id, name, type, description, location, specialties)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	err := r.pool.QueryRow(ctx, query, business.UserID, business.Name, business.Type, business.Description, business.Location, business.Specialties).Scan(&business.ID)
	if err != nil {
		return entities.Business{}, fmt.Errorf("no se pudo crear el negocio: %w", err)
	}
	return r.FindByID(ctx, business.ID)
}

func (r *PostgreSQLBusinessRepository) ExistsByUserID(ctx context.Context, userID string) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM businesses WHERE user_id = $1)`, userID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("no se pudo verificar el negocio: %w", err)
	}
	return exists, nil
}

func (r *PostgreSQLBusinessRepository) FindAll(ctx context.Context) ([]entities.Business, error) {
	rows, err := r.pool.Query(ctx, selectBusinessesQuery+" ORDER BY created_at DESC")
	if err != nil {
		return nil, fmt.Errorf("no se pudieron listar los negocios: %w", err)
	}
	defer rows.Close()

	var list []entities.Business
	for rows.Next() {
		b, err := scanBusiness(rows)
		if err != nil {
			return nil, fmt.Errorf("no se pudo leer el negocio: %w", err)
		}
		list = append(list, b)
	}
	return list, rows.Err()
}

func (r *PostgreSQLBusinessRepository) FindByID(ctx context.Context, id string) (entities.Business, error) {
	row := r.pool.QueryRow(ctx, selectBusinessesQuery+" WHERE id = $1", id)
	b, err := scanBusiness(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return entities.Business{}, repositories.ErrBusinessNotFound
	}
	if err != nil {
		return entities.Business{}, fmt.Errorf("no se pudo obtener el negocio: %w", err)
	}
	return b, nil
}

func (r *PostgreSQLBusinessRepository) Update(ctx context.Context, id string, userID string, business entities.Business) (entities.Business, error) {
	const query = `
		UPDATE businesses
		SET name = $1, type = $2, description = $3, location = $4, specialties = $5
		WHERE id = $6 AND user_id = $7
	`
	tag, err := r.pool.Exec(ctx, query, business.Name, business.Type, business.Description, business.Location, business.Specialties, id, userID)
	if err != nil {
		return entities.Business{}, fmt.Errorf("no se pudo actualizar el negocio: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return entities.Business{}, repositories.ErrBusinessNotFound
	}
	return r.FindByID(ctx, id)
}

func (r *PostgreSQLBusinessRepository) Delete(ctx context.Context, id string, userID string) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM businesses WHERE id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		return fmt.Errorf("no se pudo eliminar el negocio: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return repositories.ErrBusinessNotFound
	}
	return nil
}
