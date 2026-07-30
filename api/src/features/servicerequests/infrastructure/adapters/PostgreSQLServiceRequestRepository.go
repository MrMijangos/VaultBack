package adapters

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"vault/src/features/servicerequests/domain/entities"
)

// selectServiceRequestsQuery trae, junto con cada solicitud, el nombre e
// imagen del activo, el nombre del dueño y el nombre del negocio -- todo
// con joins (misma base de datos) para que listar no dispare una consulta
// extra por fila, mismo criterio que selectAssetsQuery en assets/.
const selectServiceRequestsQuery = `
	SELECT sr.id, sr.asset_id, sr.owner_id, sr.business_id, sr.type, sr.status,
	       sr.created_at, sr.accepted_at, sr.started_at, sr.finished_at, sr.confirmed_at,
	       a.name, COALESCE((
	           SELECT url FROM asset_photos WHERE asset_id = a.id AND is_cover = true LIMIT 1
	       ), (
	           SELECT url FROM asset_photos WHERE asset_id = a.id ORDER BY "order" LIMIT 1
	       ), ''),
	       u.name, b.name
	FROM service_requests sr
	JOIN assets a ON a.id = sr.asset_id
	JOIN users u ON u.id = sr.owner_id
	JOIN businesses b ON b.id = sr.business_id
`

type PostgreSQLServiceRequestRepository struct {
	pool *pgxpool.Pool
}

func NewPostgreSQLServiceRequestRepository(pool *pgxpool.Pool) *PostgreSQLServiceRequestRepository {
	return &PostgreSQLServiceRequestRepository{pool: pool}
}

func scanServiceRequest(row pgx.Row) (entities.ServiceRequest, error) {
	var sr entities.ServiceRequest
	err := row.Scan(
		&sr.ID, &sr.AssetID, &sr.OwnerID, &sr.BusinessID, &sr.Type, &sr.Status,
		&sr.CreatedAt, &sr.AcceptedAt, &sr.StartedAt, &sr.FinishedAt, &sr.ConfirmedAt,
		&sr.AssetName, &sr.AssetImageURL,
		&sr.OwnerName, &sr.BusinessName,
	)
	return sr, err
}

func (r *PostgreSQLServiceRequestRepository) Create(ctx context.Context, sr *entities.ServiceRequest) error {
	const query = `
		INSERT INTO service_requests (id, asset_id, owner_id, business_id, type, status)
		VALUES (gen_random_uuid(), $1, $2, $3, $4, $5)
		RETURNING id, created_at
	`
	err := r.pool.QueryRow(ctx, query, sr.AssetID, sr.OwnerID, sr.BusinessID, sr.Type, sr.Status).
		Scan(&sr.ID, &sr.CreatedAt)
	if err != nil {
		return fmt.Errorf("no se pudo crear la solicitud de servicio: %w", err)
	}
	return nil
}

func (r *PostgreSQLServiceRequestRepository) Update(ctx context.Context, sr *entities.ServiceRequest) error {
	const query = `
		UPDATE service_requests
		SET status = $1, accepted_at = $2, started_at = $3, finished_at = $4, confirmed_at = $5
		WHERE id = $6
	`
	tag, err := r.pool.Exec(ctx, query, sr.Status, sr.AcceptedAt, sr.StartedAt, sr.FinishedAt, sr.ConfirmedAt, sr.ID)
	if err != nil {
		return fmt.Errorf("no se pudo actualizar la solicitud de servicio: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("solicitud de servicio %q no existe", sr.ID)
	}
	return nil
}

func (r *PostgreSQLServiceRequestRepository) FindByID(ctx context.Context, id string) (*entities.ServiceRequest, error) {
	row := r.pool.QueryRow(ctx, selectServiceRequestsQuery+" WHERE sr.id = $1", id)
	sr, err := scanServiceRequest(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("no se pudo obtener la solicitud de servicio: %w", err)
	}
	return &sr, nil
}

func (r *PostgreSQLServiceRequestRepository) ListByOwnerID(ctx context.Context, ownerID string) ([]entities.ServiceRequest, error) {
	rows, err := r.pool.Query(ctx, selectServiceRequestsQuery+" WHERE sr.owner_id = $1 ORDER BY sr.created_at DESC", ownerID)
	if err != nil {
		return nil, fmt.Errorf("no se pudieron listar las solicitudes de servicio: %w", err)
	}
	defer rows.Close()

	var out []entities.ServiceRequest
	for rows.Next() {
		sr, err := scanServiceRequest(rows)
		if err != nil {
			return nil, fmt.Errorf("no se pudo leer la solicitud de servicio: %w", err)
		}
		out = append(out, sr)
	}
	return out, rows.Err()
}

func (r *PostgreSQLServiceRequestRepository) ListByBusinessID(ctx context.Context, businessID string) ([]entities.ServiceRequest, error) {
	rows, err := r.pool.Query(ctx, selectServiceRequestsQuery+" WHERE sr.business_id = $1 ORDER BY sr.created_at DESC", businessID)
	if err != nil {
		return nil, fmt.Errorf("no se pudieron listar las solicitudes de servicio: %w", err)
	}
	defer rows.Close()

	var out []entities.ServiceRequest
	for rows.Next() {
		sr, err := scanServiceRequest(rows)
		if err != nil {
			return nil, fmt.Errorf("no se pudo leer la solicitud de servicio: %w", err)
		}
		out = append(out, sr)
	}
	return out, rows.Err()
}
