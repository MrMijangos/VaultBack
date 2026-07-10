package adapters

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"vault/src/features/reviews/domain/entities"
	"vault/src/features/reviews/domain/repositories"
)

const selectReviewsQuery = `
	SELECT id, user_id, provider_id, content, is_visible, likes_count, created_at
	FROM reviews
`

type PostgreSQLReviewRepository struct {
	pool *pgxpool.Pool
}

func NewPostgreSQLReviewRepository(pool *pgxpool.Pool) *PostgreSQLReviewRepository {
	return &PostgreSQLReviewRepository{pool: pool}
}

func scanReview(row pgx.Row) (entities.Review, error) {
	var r entities.Review
	err := row.Scan(&r.ID, &r.UserID, &r.ProviderID, &r.Content, &r.IsVisible, &r.LikesCount, &r.CreatedAt)
	return r, err
}

func (r *PostgreSQLReviewRepository) Create(ctx context.Context, review entities.Review) (entities.Review, error) {
	const query = `
		INSERT INTO reviews (user_id, provider_id, content)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	err := r.pool.QueryRow(ctx, query, review.UserID, review.ProviderID, review.Content).Scan(&review.ID)
	if err != nil {
		return entities.Review{}, fmt.Errorf("no se pudo crear la reseña: %w", err)
	}
	return r.FindByID(ctx, review.ID)
}

func (r *PostgreSQLReviewRepository) FindByProviderID(ctx context.Context, providerID string) ([]entities.Review, error) {
	rows, err := r.pool.Query(ctx, selectReviewsQuery+" WHERE provider_id = $1 AND is_visible = true ORDER BY created_at DESC", providerID)
	if err != nil {
		return nil, fmt.Errorf("no se pudieron listar las reseñas: %w", err)
	}
	defer rows.Close()

	var list []entities.Review
	for rows.Next() {
		rv, err := scanReview(rows)
		if err != nil {
			return nil, fmt.Errorf("no se pudo leer la reseña: %w", err)
		}
		list = append(list, rv)
	}
	return list, rows.Err()
}

func (r *PostgreSQLReviewRepository) FindByID(ctx context.Context, id string) (entities.Review, error) {
	row := r.pool.QueryRow(ctx, selectReviewsQuery+" WHERE id = $1", id)
	rv, err := scanReview(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return entities.Review{}, repositories.ErrReviewNotFound
	}
	if err != nil {
		return entities.Review{}, fmt.Errorf("no se pudo obtener la reseña: %w", err)
	}
	return rv, nil
}

func (r *PostgreSQLReviewRepository) Delete(ctx context.Context, id string, userID string) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM reviews WHERE id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		return fmt.Errorf("no se pudo eliminar la reseña: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return repositories.ErrReviewNotFound
	}
	return nil
}

func (r *PostgreSQLReviewRepository) Like(ctx context.Context, reviewID string, userID string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("no se pudo iniciar la transaccion: %w", err)
	}
	defer tx.Rollback(ctx)

	tag, err := tx.Exec(ctx, `INSERT INTO review_likes (review_id, user_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`, reviewID, userID)
	if err != nil {
		return fmt.Errorf("no se pudo dar like: %w", err)
	}
	if tag.RowsAffected() > 0 {
		if _, err := tx.Exec(ctx, `UPDATE reviews SET likes_count = likes_count + 1 WHERE id = $1`, reviewID); err != nil {
			return fmt.Errorf("no se pudo actualizar el contador de likes: %w", err)
		}
	}

	return tx.Commit(ctx)
}

func (r *PostgreSQLReviewRepository) Unlike(ctx context.Context, reviewID string, userID string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("no se pudo iniciar la transaccion: %w", err)
	}
	defer tx.Rollback(ctx)

	tag, err := tx.Exec(ctx, `DELETE FROM review_likes WHERE review_id = $1 AND user_id = $2`, reviewID, userID)
	if err != nil {
		return fmt.Errorf("no se pudo quitar el like: %w", err)
	}
	if tag.RowsAffected() > 0 {
		if _, err := tx.Exec(ctx, `UPDATE reviews SET likes_count = GREATEST(likes_count - 1, 0) WHERE id = $1`, reviewID); err != nil {
			return fmt.Errorf("no se pudo actualizar el contador de likes: %w", err)
		}
	}

	return tx.Commit(ctx)
}
