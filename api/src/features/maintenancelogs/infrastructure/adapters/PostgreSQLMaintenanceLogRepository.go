package adapters

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"vault/src/features/maintenancelogs/domain/entities"
	"vault/src/features/maintenancelogs/domain/repositories"
)

const selectMaintenanceLogsQuery = `
	SELECT id, asset_id, provider_id, type, COALESCE(subtype, ''), cost::float8,
	       performed_at, COALESCE(notes, ''), COALESCE(blockchain_tx_id, ''), created_at
	FROM maintenance_logs
`

type PostgreSQLMaintenanceLogRepository struct {
	pool *pgxpool.Pool
}

func NewPostgreSQLMaintenanceLogRepository(pool *pgxpool.Pool) *PostgreSQLMaintenanceLogRepository {
	return &PostgreSQLMaintenanceLogRepository{pool: pool}
}

func scanMaintenanceLog(row pgx.Row) (entities.MaintenanceLog, error) {
	var l entities.MaintenanceLog
	err := row.Scan(&l.ID, &l.AssetID, &l.ProviderID, &l.Type, &l.Subtype, &l.Cost, &l.PerformedAt, &l.Notes, &l.BlockchainTxID, &l.CreatedAt)
	return l, err
}

func (r *PostgreSQLMaintenanceLogRepository) Create(ctx context.Context, log entities.MaintenanceLog) (entities.MaintenanceLog, error) {
	const query = `
		INSERT INTO maintenance_logs (asset_id, provider_id, type, subtype, cost, performed_at, notes)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`
	err := r.pool.QueryRow(ctx, query, log.AssetID, log.ProviderID, log.Type, log.Subtype, log.Cost, log.PerformedAt, log.Notes).Scan(&log.ID)
	if err != nil {
		return entities.MaintenanceLog{}, fmt.Errorf("no se pudo crear el registro de servicio: %w", err)
	}
	return r.FindByID(ctx, log.ID)
}

func (r *PostgreSQLMaintenanceLogRepository) FindByAssetID(ctx context.Context, assetID string) ([]entities.MaintenanceLog, error) {
	rows, err := r.pool.Query(ctx, selectMaintenanceLogsQuery+" WHERE asset_id = $1 ORDER BY created_at DESC", assetID)
	if err != nil {
		return nil, fmt.Errorf("no se pudieron listar los registros de servicio: %w", err)
	}
	defer rows.Close()

	var list []entities.MaintenanceLog
	for rows.Next() {
		l, err := scanMaintenanceLog(rows)
		if err != nil {
			return nil, fmt.Errorf("no se pudo leer el registro de servicio: %w", err)
		}
		list = append(list, l)
	}
	return list, rows.Err()
}

func (r *PostgreSQLMaintenanceLogRepository) FindByID(ctx context.Context, id string) (entities.MaintenanceLog, error) {
	row := r.pool.QueryRow(ctx, selectMaintenanceLogsQuery+" WHERE id = $1", id)
	l, err := scanMaintenanceLog(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return entities.MaintenanceLog{}, repositories.ErrMaintenanceLogNotFound
	}
	if err != nil {
		return entities.MaintenanceLog{}, fmt.Errorf("no se pudo obtener el registro de servicio: %w", err)
	}
	return l, nil
}

func (r *PostgreSQLMaintenanceLogRepository) Update(ctx context.Context, id string, log entities.MaintenanceLog) (entities.MaintenanceLog, error) {
	const query = `
		UPDATE maintenance_logs
		SET type = $1, subtype = $2, cost = $3, performed_at = $4, notes = $5
		WHERE id = $6
	`
	tag, err := r.pool.Exec(ctx, query, log.Type, log.Subtype, log.Cost, log.PerformedAt, log.Notes, id)
	if err != nil {
		return entities.MaintenanceLog{}, fmt.Errorf("no se pudo actualizar el registro de servicio: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return entities.MaintenanceLog{}, repositories.ErrMaintenanceLogNotFound
	}
	return r.FindByID(ctx, id)
}

func (r *PostgreSQLMaintenanceLogRepository) Delete(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM maintenance_logs WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("no se pudo eliminar el registro de servicio: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return repositories.ErrMaintenanceLogNotFound
	}
	return nil
}
