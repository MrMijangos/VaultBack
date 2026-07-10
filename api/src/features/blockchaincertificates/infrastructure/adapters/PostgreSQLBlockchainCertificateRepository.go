package adapters

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"vault/src/features/blockchaincertificates/domain/entities"
	"vault/src/features/blockchaincertificates/domain/repositories"
)

const selectCertificatesQuery = `
	SELECT id, asset_id, owner_id, tx_id, asset_hash, action, network, confirmed_at
	FROM blockchain_certificates
`

type PostgreSQLBlockchainCertificateRepository struct {
	pool *pgxpool.Pool
}

func NewPostgreSQLBlockchainCertificateRepository(pool *pgxpool.Pool) *PostgreSQLBlockchainCertificateRepository {
	return &PostgreSQLBlockchainCertificateRepository{pool: pool}
}

func scanCertificate(row pgx.Row) (entities.BlockchainCertificate, error) {
	var c entities.BlockchainCertificate
	err := row.Scan(&c.ID, &c.AssetID, &c.OwnerID, &c.TxID, &c.AssetHash, &c.Action, &c.Network, &c.ConfirmedAt)
	return c, err
}

func (r *PostgreSQLBlockchainCertificateRepository) Create(ctx context.Context, cert entities.BlockchainCertificate) (entities.BlockchainCertificate, error) {
	const query = `
		INSERT INTO blockchain_certificates (asset_id, owner_id, tx_id, asset_hash, action, network)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	err := r.pool.QueryRow(ctx, query, cert.AssetID, cert.OwnerID, cert.TxID, cert.AssetHash, cert.Action, cert.Network).Scan(&cert.ID)
	if err != nil {
		return entities.BlockchainCertificate{}, fmt.Errorf("no se pudo crear el certificado: %w", err)
	}
	return r.FindByID(ctx, cert.ID)
}

func (r *PostgreSQLBlockchainCertificateRepository) ExistsByTxID(ctx context.Context, txID string) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM blockchain_certificates WHERE tx_id = $1)`, txID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("no se pudo verificar el tx_id: %w", err)
	}
	return exists, nil
}

func (r *PostgreSQLBlockchainCertificateRepository) FindByAssetID(ctx context.Context, assetID string) ([]entities.BlockchainCertificate, error) {
	rows, err := r.pool.Query(ctx, selectCertificatesQuery+" WHERE asset_id = $1 ORDER BY confirmed_at DESC", assetID)
	if err != nil {
		return nil, fmt.Errorf("no se pudieron listar los certificados: %w", err)
	}
	defer rows.Close()

	var list []entities.BlockchainCertificate
	for rows.Next() {
		c, err := scanCertificate(rows)
		if err != nil {
			return nil, fmt.Errorf("no se pudo leer el certificado: %w", err)
		}
		list = append(list, c)
	}
	return list, rows.Err()
}

func (r *PostgreSQLBlockchainCertificateRepository) FindByID(ctx context.Context, id string) (entities.BlockchainCertificate, error) {
	row := r.pool.QueryRow(ctx, selectCertificatesQuery+" WHERE id = $1", id)
	c, err := scanCertificate(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return entities.BlockchainCertificate{}, repositories.ErrCertificateNotFound
	}
	if err != nil {
		return entities.BlockchainCertificate{}, fmt.Errorf("no se pudo obtener el certificado: %w", err)
	}
	return c, nil
}
