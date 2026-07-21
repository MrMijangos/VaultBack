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
	SELECT r.id, r.user_id, r.provider_id, r.content, r.sentiment_score, r.sentiment_label,
	       r.is_visible, r.likes_count, r.created_at,
	       u.name, COALESCE(u.avatar_url, '')
	FROM reviews r
	JOIN users u ON u.id = r.user_id
`

// providerRatingQuery convierte cada reseña visible y ya analizada en una
// puntuación de 1 a 10: 5.5 es el punto neutral, y se desplaza hacia 10
// (positivo) o hacia 1 (negativo) según la confianza del modelo en esa
// etiqueta -- ver AnalyzeSentiment en el servicio de NLP. Promediar esto en
// vez del sentiment_score crudo es necesario porque ese score es solo la
// confianza del modelo en su etiqueta (siempre alto sin importar si es
// positiva o negativa), no una medida de qué tan buena es la reseña.
const providerRatingQuery = `
	SELECT
		AVG(
			CASE sentiment_label
				WHEN 'positivo' THEN 5.5 + sentiment_score * 4.5
				WHEN 'negativo' THEN 5.5 - sentiment_score * 4.5
				ELSE 5.5
			END
		),
		COUNT(*)
	FROM reviews
	WHERE provider_id = $1 AND is_visible = true AND sentiment_label IS NOT NULL
`

type PostgreSQLReviewRepository struct {
	pool *pgxpool.Pool
}

func NewPostgreSQLReviewRepository(pool *pgxpool.Pool) *PostgreSQLReviewRepository {
	return &PostgreSQLReviewRepository{pool: pool}
}

func scanReview(row pgx.Row) (entities.Review, error) {
	var r entities.Review
	err := row.Scan(
		&r.ID, &r.UserID, &r.ProviderID, &r.Content, &r.SentimentScore, &r.SentimentLabel,
		&r.IsVisible, &r.LikesCount, &r.CreatedAt,
		&r.AuthorName, &r.AuthorAvatarURL,
	)
	return r, err
}

func (r *PostgreSQLReviewRepository) Create(ctx context.Context, review entities.Review) (entities.Review, error) {
	const query = `
		INSERT INTO reviews (id, user_id, provider_id, content, sentiment_score, toxicity_score, is_visible)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`
	err := r.pool.QueryRow(ctx, query,
		review.ID, review.UserID, review.ProviderID, review.Content,
		review.SentimentScore, review.ToxicityScore, review.IsVisible,
	).Scan(&review.ID)
	if err != nil {
		return entities.Review{}, fmt.Errorf("no se pudo crear la reseña: %w", err)
	}
	return r.FindByID(ctx, review.ID)
}

func (r *PostgreSQLReviewRepository) FindByProviderID(ctx context.Context, providerID string) ([]entities.Review, error) {
	rows, err := r.pool.Query(ctx, selectReviewsQuery+" WHERE r.provider_id = $1 AND r.is_visible = true ORDER BY r.created_at DESC", providerID)
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
	row := r.pool.QueryRow(ctx, selectReviewsQuery+" WHERE r.id = $1", id)
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

func (r *PostgreSQLReviewRepository) GetProviderRating(ctx context.Context, providerID string) (*float64, int, error) {
	var rating *float64
	var total int
	err := r.pool.QueryRow(ctx, providerRatingQuery, providerID).Scan(&rating, &total)
	if err != nil {
		return nil, 0, fmt.Errorf("no se pudo calcular la calificación del proveedor: %w", err)
	}
	return rating, total, nil
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
