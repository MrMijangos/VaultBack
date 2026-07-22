package adapters

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"vault-payment/src/features/ads/domain/entities"
)

const selectAdsQuery = `
	SELECT id, user_id, subscription_id, title, COALESCE(description, ''), COALESCE(image_url, ''),
	       target_section, COALESCE(target_id, ''), status, impressions, clicks, created_at
	FROM ads
`

// PostgreSQLAdRepository reemplaza InMemoryAdRepository.
type PostgreSQLAdRepository struct {
	pool *pgxpool.Pool
}

func NewPostgreSQLAdRepository(pool *pgxpool.Pool) *PostgreSQLAdRepository {
	return &PostgreSQLAdRepository{pool: pool}
}

func scanAd(row pgx.Row) (*entities.Ad, error) {
	var a entities.Ad
	err := row.Scan(
		&a.ID, &a.UserID, &a.SubscriptionID, &a.Title, &a.Description, &a.ImageURL,
		&a.TargetSection, &a.TargetID, &a.Status, &a.Impressions, &a.Clicks, &a.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func scanAds(rows pgx.Rows) ([]*entities.Ad, error) {
	defer rows.Close()
	var out []*entities.Ad
	for rows.Next() {
		a, err := scanAd(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

func (r *PostgreSQLAdRepository) Create(ctx context.Context, ad *entities.Ad) error {
	const query = `
		INSERT INTO ads (id, user_id, subscription_id, title, description, image_url, target_section, target_id, status, impressions, clicks, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`
	_, err := r.pool.Exec(ctx, query,
		ad.ID, ad.UserID, ad.SubscriptionID, ad.Title, ad.Description, ad.ImageURL,
		ad.TargetSection, ad.TargetID, ad.Status, ad.Impressions, ad.Clicks, ad.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("no se pudo crear el anuncio: %w", err)
	}
	return nil
}

func (r *PostgreSQLAdRepository) Update(ctx context.Context, ad *entities.Ad) error {
	const query = `
		UPDATE ads
		SET title = $1, description = $2, image_url = $3, target_section = $4, target_id = $5, status = $6
		WHERE id = $7
	`
	tag, err := r.pool.Exec(ctx, query, ad.Title, ad.Description, ad.ImageURL, ad.TargetSection, ad.TargetID, ad.Status, ad.ID)
	if err != nil {
		return fmt.Errorf("no se pudo actualizar el anuncio: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("anuncio %q no existe", ad.ID)
	}
	return nil
}

func (r *PostgreSQLAdRepository) Delete(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM ads WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("no se pudo eliminar el anuncio: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("anuncio %q no existe", id)
	}
	return nil
}

func (r *PostgreSQLAdRepository) GetByID(ctx context.Context, id string) (*entities.Ad, error) {
	row := r.pool.QueryRow(ctx, selectAdsQuery+" WHERE id = $1", id)
	ad, err := scanAd(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("no se pudo obtener el anuncio: %w", err)
	}
	return ad, nil
}

func (r *PostgreSQLAdRepository) ListByUserID(ctx context.Context, userID string) ([]*entities.Ad, error) {
	rows, err := r.pool.Query(ctx, selectAdsQuery+" WHERE user_id = $1 ORDER BY created_at DESC", userID)
	if err != nil {
		return nil, fmt.Errorf("no se pudieron listar los anuncios: %w", err)
	}
	return scanAds(rows)
}

func (r *PostgreSQLAdRepository) ListActiveBySection(ctx context.Context, section string) ([]*entities.Ad, error) {
	rows, err := r.pool.Query(ctx, selectAdsQuery+" WHERE status = $1 AND target_section = $2 ORDER BY created_at DESC", entities.AdStatusActive, section)
	if err != nil {
		return nil, fmt.Errorf("no se pudieron listar los anuncios: %w", err)
	}
	return scanAds(rows)
}

func (r *PostgreSQLAdRepository) ListBySubscriptionID(ctx context.Context, subscriptionID string) ([]*entities.Ad, error) {
	rows, err := r.pool.Query(ctx, selectAdsQuery+" WHERE subscription_id = $1", subscriptionID)
	if err != nil {
		return nil, fmt.Errorf("no se pudieron listar los anuncios: %w", err)
	}
	return scanAds(rows)
}

func (r *PostgreSQLAdRepository) IncrementImpressions(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, `UPDATE ads SET impressions = impressions + 1 WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("no se pudo registrar la impresión: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("anuncio %q no existe", id)
	}
	return nil
}

func (r *PostgreSQLAdRepository) IncrementClicks(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, `UPDATE ads SET clicks = clicks + 1 WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("no se pudo registrar el clic: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("anuncio %q no existe", id)
	}
	return nil
}

func (r *PostgreSQLAdRepository) CountActiveByUserID(ctx context.Context, userID string) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM ads WHERE user_id = $1 AND status = $2`,
		userID, entities.AdStatusActive,
	).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("no se pudo contar los anuncios activos: %w", err)
	}
	return count, nil
}
