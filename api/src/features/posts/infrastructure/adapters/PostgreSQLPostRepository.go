package adapters

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"vault/src/features/posts/domain/entities"
	"vault/src/features/posts/domain/repositories"
)

const selectPostsQuery = `
	SELECT id, user_id, asset_id, content, sentiment_score, COALESCE(sentiment_label, ''),
	       toxicity_score, is_visible, likes_count, created_at
	FROM posts
`

type PostgreSQLPostRepository struct {
	pool *pgxpool.Pool
}

func NewPostgreSQLPostRepository(pool *pgxpool.Pool) *PostgreSQLPostRepository {
	return &PostgreSQLPostRepository{pool: pool}
}

func scanPost(row pgx.Row) (entities.Post, error) {
	var p entities.Post
	err := row.Scan(&p.ID, &p.UserID, &p.AssetID, &p.Content, &p.SentimentScore, &p.SentimentLabel, &p.ToxicityScore, &p.IsVisible, &p.LikesCount, &p.CreatedAt)
	return p, err
}

func (r *PostgreSQLPostRepository) Create(ctx context.Context, post entities.Post) (entities.Post, error) {
	const query = `
		INSERT INTO posts (user_id, asset_id, content)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	err := r.pool.QueryRow(ctx, query, post.UserID, post.AssetID, post.Content).Scan(&post.ID)
	if err != nil {
		return entities.Post{}, fmt.Errorf("no se pudo crear la publicacion: %w", err)
	}
	return r.FindByID(ctx, post.ID)
}

func (r *PostgreSQLPostRepository) FindAllVisible(ctx context.Context) ([]entities.Post, error) {
	rows, err := r.pool.Query(ctx, selectPostsQuery+" WHERE is_visible = true ORDER BY created_at DESC")
	if err != nil {
		return nil, fmt.Errorf("no se pudieron listar las publicaciones: %w", err)
	}
	defer rows.Close()

	var list []entities.Post
	for rows.Next() {
		p, err := scanPost(rows)
		if err != nil {
			return nil, fmt.Errorf("no se pudo leer la publicacion: %w", err)
		}
		list = append(list, p)
	}
	return list, rows.Err()
}

func (r *PostgreSQLPostRepository) FindByID(ctx context.Context, id string) (entities.Post, error) {
	row := r.pool.QueryRow(ctx, selectPostsQuery+" WHERE id = $1", id)
	p, err := scanPost(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return entities.Post{}, repositories.ErrPostNotFound
	}
	if err != nil {
		return entities.Post{}, fmt.Errorf("no se pudo obtener la publicacion: %w", err)
	}
	return p, nil
}

func (r *PostgreSQLPostRepository) Update(ctx context.Context, id string, userID string, content string) (entities.Post, error) {
	tag, err := r.pool.Exec(ctx, `UPDATE posts SET content = $1 WHERE id = $2 AND user_id = $3`, content, id, userID)
	if err != nil {
		return entities.Post{}, fmt.Errorf("no se pudo actualizar la publicacion: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return entities.Post{}, repositories.ErrPostNotFound
	}
	return r.FindByID(ctx, id)
}

func (r *PostgreSQLPostRepository) Delete(ctx context.Context, id string, userID string) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM posts WHERE id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		return fmt.Errorf("no se pudo eliminar la publicacion: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return repositories.ErrPostNotFound
	}
	return nil
}

func (r *PostgreSQLPostRepository) FindPhotosByPostID(ctx context.Context, postID string) ([]entities.PostPhoto, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, post_id, url, "order", created_at FROM post_photos WHERE post_id = $1 ORDER BY "order"`, postID)
	if err != nil {
		return nil, fmt.Errorf("no se pudieron listar las fotos: %w", err)
	}
	defer rows.Close()

	var photos []entities.PostPhoto
	for rows.Next() {
		var p entities.PostPhoto
		if err := rows.Scan(&p.ID, &p.PostID, &p.URL, &p.Order, &p.CreatedAt); err != nil {
			return nil, fmt.Errorf("no se pudo leer la foto: %w", err)
		}
		photos = append(photos, p)
	}
	return photos, rows.Err()
}

func (r *PostgreSQLPostRepository) AddPhoto(ctx context.Context, postID string, url string) (entities.PostPhoto, error) {
	var count int
	if err := r.pool.QueryRow(ctx, `SELECT COUNT(*) FROM post_photos WHERE post_id = $1`, postID).Scan(&count); err != nil {
		return entities.PostPhoto{}, fmt.Errorf("no se pudo verificar las fotos existentes: %w", err)
	}

	const query = `
		INSERT INTO post_photos (post_id, url, "order")
		VALUES ($1, $2, $3)
		RETURNING id, post_id, url, "order", created_at
	`
	var p entities.PostPhoto
	err := r.pool.QueryRow(ctx, query, postID, url, count).Scan(&p.ID, &p.PostID, &p.URL, &p.Order, &p.CreatedAt)
	if err != nil {
		return entities.PostPhoto{}, fmt.Errorf("no se pudo guardar la foto: %w", err)
	}
	return p, nil
}

func (r *PostgreSQLPostRepository) Like(ctx context.Context, postID string, userID string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("no se pudo iniciar la transaccion: %w", err)
	}
	defer tx.Rollback(ctx)

	tag, err := tx.Exec(ctx, `INSERT INTO post_likes (post_id, user_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`, postID, userID)
	if err != nil {
		return fmt.Errorf("no se pudo dar like: %w", err)
	}
	if tag.RowsAffected() > 0 {
		if _, err := tx.Exec(ctx, `UPDATE posts SET likes_count = likes_count + 1 WHERE id = $1`, postID); err != nil {
			return fmt.Errorf("no se pudo actualizar el contador de likes: %w", err)
		}
	}

	return tx.Commit(ctx)
}

func (r *PostgreSQLPostRepository) Unlike(ctx context.Context, postID string, userID string) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("no se pudo iniciar la transaccion: %w", err)
	}
	defer tx.Rollback(ctx)

	tag, err := tx.Exec(ctx, `DELETE FROM post_likes WHERE post_id = $1 AND user_id = $2`, postID, userID)
	if err != nil {
		return fmt.Errorf("no se pudo quitar el like: %w", err)
	}
	if tag.RowsAffected() > 0 {
		if _, err := tx.Exec(ctx, `UPDATE posts SET likes_count = GREATEST(likes_count - 1, 0) WHERE id = $1`, postID); err != nil {
			return fmt.Errorf("no se pudo actualizar el contador de likes: %w", err)
		}
	}

	return tx.Commit(ctx)
}
