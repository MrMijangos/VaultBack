package adapters

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"vault/src/features/businessservices/domain/entities"
	"vault/src/features/businessservices/domain/repositories"
)

const selectBusinessServicesQuery = `
	SELECT id, business_id, title, COALESCE(description, ''), price::float8, created_at
	FROM business_services
`

type PostgreSQLBusinessServiceRepository struct {
	pool *pgxpool.Pool
}

func NewPostgreSQLBusinessServiceRepository(pool *pgxpool.Pool) *PostgreSQLBusinessServiceRepository {
	return &PostgreSQLBusinessServiceRepository{pool: pool}
}

func scanBusinessService(row pgx.Row) (entities.BusinessService, error) {
	var s entities.BusinessService
	err := row.Scan(&s.ID, &s.BusinessID, &s.Title, &s.Description, &s.Price, &s.CreatedAt)
	return s, err
}

func (r *PostgreSQLBusinessServiceRepository) Create(ctx context.Context, service entities.BusinessService) (entities.BusinessService, error) {
	const query = `
		INSERT INTO business_services (business_id, title, description, price)
		VALUES ($1, $2, $3, $4)
		RETURNING id, business_id, title, COALESCE(description, ''), price::float8, created_at
	`
	row := r.pool.QueryRow(ctx, query, service.BusinessID, service.Title, service.Description, service.Price)
	created, err := scanBusinessService(row)
	if err != nil {
		return entities.BusinessService{}, fmt.Errorf("no se pudo crear el servicio: %w", err)
	}
	return created, nil
}

func (r *PostgreSQLBusinessServiceRepository) Update(ctx context.Context, id string, businessID string, service entities.BusinessService) (entities.BusinessService, error) {
	const query = `
		UPDATE business_services
		SET title = $1, description = $2, price = $3
		WHERE id = $4 AND business_id = $5
	`
	tag, err := r.pool.Exec(ctx, query, service.Title, service.Description, service.Price, id, businessID)
	if err != nil {
		return entities.BusinessService{}, fmt.Errorf("no se pudo actualizar el servicio: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return entities.BusinessService{}, repositories.ErrBusinessServiceNotFound
	}

	row := r.pool.QueryRow(ctx, selectBusinessServicesQuery+" WHERE id = $1", id)
	updated, err := scanBusinessService(row)
	if err != nil {
		return entities.BusinessService{}, fmt.Errorf("no se pudo leer el servicio actualizado: %w", err)
	}
	return updated, nil
}

func (r *PostgreSQLBusinessServiceRepository) Delete(ctx context.Context, id string, businessID string) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM business_services WHERE id = $1 AND business_id = $2`, id, businessID)
	if err != nil {
		return fmt.Errorf("no se pudo eliminar el servicio: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return repositories.ErrBusinessServiceNotFound
	}
	return nil
}

func (r *PostgreSQLBusinessServiceRepository) ListByBusinessID(ctx context.Context, businessID string) ([]entities.BusinessService, error) {
	rows, err := r.pool.Query(ctx, selectBusinessServicesQuery+" WHERE business_id = $1 ORDER BY created_at", businessID)
	if err != nil {
		return nil, fmt.Errorf("no se pudieron listar los servicios: %w", err)
	}
	defer rows.Close()

	var list []entities.BusinessService
	for rows.Next() {
		s, err := scanBusinessService(rows)
		if err != nil {
			return nil, fmt.Errorf("no se pudo leer el servicio: %w", err)
		}
		list = append(list, s)
	}
	return list, rows.Err()
}
