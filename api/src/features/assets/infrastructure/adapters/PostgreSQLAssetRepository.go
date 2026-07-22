package adapters

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"vault/src/features/assets/domain/entities"
	"vault/src/features/assets/domain/repositories"
)

const selectAssetsQuery = `
	SELECT id, user_id, name, category, COALESCE(brand, ''), purchase_value::float8, condition,
	       purchase_date, COALESCE(store_origin, ''), COALESCE(notes, ''),
	       COALESCE(blockchain_tx_id, ''), COALESCE(blockchain_hash, ''), created_at,
	       COALESCE(is_for_sale, false), sale_price::float8, COALESCE(sale_description, ''), COALESCE(size, '')
	FROM assets
`

type PostgreSQLAssetRepository struct {
	pool *pgxpool.Pool
}

func NewPostgreSQLAssetRepository(pool *pgxpool.Pool) *PostgreSQLAssetRepository {
	return &PostgreSQLAssetRepository{pool: pool}
}

func scanAsset(row pgx.Row) (entities.Asset, error) {
	var a entities.Asset
	err := row.Scan(
		&a.ID, &a.UserID, &a.Name, &a.Category, &a.Brand, &a.PurchaseValue, &a.Condition,
		&a.PurchaseDate, &a.StoreOrigin, &a.Notes, &a.BlockchainTxID, &a.BlockchainHash, &a.CreatedAt,
		&a.IsForSale, &a.SalePrice, &a.SaleDescription, &a.Size,
	)
	return a, err
}

func (r *PostgreSQLAssetRepository) Create(ctx context.Context, asset entities.Asset) (entities.Asset, error) {
	const query = `
		INSERT INTO assets (user_id, name, category, brand, purchase_value, condition, purchase_date, store_origin, notes, is_for_sale, sale_price, sale_description, size)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id
	`

	err := r.pool.QueryRow(ctx, query,
		asset.UserID, asset.Name, asset.Category, asset.Brand, asset.PurchaseValue,
		asset.Condition, asset.PurchaseDate, asset.StoreOrigin, asset.Notes,
		asset.IsForSale, asset.SalePrice, asset.SaleDescription, asset.Size,
	).Scan(&asset.ID)
	if err != nil {
		return entities.Asset{}, fmt.Errorf("no se pudo crear el producto: %w", err)
	}

	return r.FindByID(ctx, asset.ID)
}

func (r *PostgreSQLAssetRepository) FindAll(ctx context.Context) ([]entities.Asset, error) {
	rows, err := r.pool.Query(ctx, selectAssetsQuery+" ORDER BY created_at DESC")
	if err != nil {
		return nil, fmt.Errorf("no se pudieron listar los productos: %w", err)
	}
	defer rows.Close()

	var assets []entities.Asset
	for rows.Next() {
		a, err := scanAsset(rows)
		if err != nil {
			return nil, fmt.Errorf("no se pudo leer el producto: %w", err)
		}
		assets = append(assets, a)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error al listar productos: %w", err)
	}

	return assets, nil
}

func (r *PostgreSQLAssetRepository) FindByID(ctx context.Context, id string) (entities.Asset, error) {
	row := r.pool.QueryRow(ctx, selectAssetsQuery+" WHERE id = $1", id)
	a, err := scanAsset(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return entities.Asset{}, repositories.ErrAssetNotFound
	}
	if err != nil {
		return entities.Asset{}, fmt.Errorf("no se pudo obtener el producto: %w", err)
	}

	return a, nil
}

func (r *PostgreSQLAssetRepository) Update(ctx context.Context, id string, userID string, asset entities.Asset) (entities.Asset, error) {
	const query = `
		UPDATE assets
		SET name = $1, category = $2, brand = $3, purchase_value = $4,
		    condition = $5, purchase_date = $6, store_origin = $7, notes = $8,
		    is_for_sale = $9, sale_price = $10, sale_description = $11, size = $12
		WHERE id = $13 AND user_id = $14
	`

	tag, err := r.pool.Exec(ctx, query,
		asset.Name, asset.Category, asset.Brand, asset.PurchaseValue,
		asset.Condition, asset.PurchaseDate, asset.StoreOrigin, asset.Notes,
		asset.IsForSale, asset.SalePrice, asset.SaleDescription, asset.Size,
		id, userID,
	)
	if err != nil {
		return entities.Asset{}, fmt.Errorf("no se pudo actualizar el producto: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return entities.Asset{}, repositories.ErrAssetNotFound
	}

	return r.FindByID(ctx, id)
}

func (r *PostgreSQLAssetRepository) Delete(ctx context.Context, id string, userID string) error {
	const query = `DELETE FROM assets WHERE id = $1 AND user_id = $2`

	tag, err := r.pool.Exec(ctx, query, id, userID)
	if err != nil {
		return fmt.Errorf("no se pudo eliminar el producto: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return repositories.ErrAssetNotFound
	}

	return nil
}

func (r *PostgreSQLAssetRepository) FindPhotosByAssetID(ctx context.Context, assetID string) ([]entities.AssetPhoto, error) {
	const query = `
		SELECT id, asset_id, url, is_cover, "order", created_at
		FROM asset_photos
		WHERE asset_id = $1
		ORDER BY "order"
	`

	rows, err := r.pool.Query(ctx, query, assetID)
	if err != nil {
		return nil, fmt.Errorf("no se pudieron listar las fotos: %w", err)
	}
	defer rows.Close()

	var photos []entities.AssetPhoto
	for rows.Next() {
		var p entities.AssetPhoto
		if err := rows.Scan(&p.ID, &p.AssetID, &p.URL, &p.IsCover, &p.Order, &p.CreatedAt); err != nil {
			return nil, fmt.Errorf("no se pudo leer la foto: %w", err)
		}
		photos = append(photos, p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error al listar fotos: %w", err)
	}

	return photos, nil
}

func (r *PostgreSQLAssetRepository) AddPhoto(ctx context.Context, assetID string, url string) (entities.AssetPhoto, error) {
	var count int
	if err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM asset_photos WHERE asset_id = $1`, assetID).Scan(&count); err != nil {
		return entities.AssetPhoto{}, fmt.Errorf("no se pudo verificar las fotos existentes: %w", err)
	}

	const query = `
		INSERT INTO asset_photos (asset_id, url, is_cover, "order")
		VALUES ($1, $2, $3, $4)
		RETURNING id, asset_id, url, is_cover, "order", created_at
	`

	var p entities.AssetPhoto
	err := r.pool.QueryRow(ctx, query, assetID, url, count == 0, count).Scan(
		&p.ID, &p.AssetID, &p.URL, &p.IsCover, &p.Order, &p.CreatedAt,
	)
	if err != nil {
		return entities.AssetPhoto{}, fmt.Errorf("no se pudo guardar la foto: %w", err)
	}

	return p, nil
}
