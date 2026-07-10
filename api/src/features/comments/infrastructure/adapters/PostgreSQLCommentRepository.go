package adapters

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"vault/src/features/comments/domain/entities"
	"vault/src/features/comments/domain/repositories"
)

type PostgreSQLCommentRepository struct {
	pool *pgxpool.Pool
}

func NewPostgreSQLCommentRepository(pool *pgxpool.Pool) *PostgreSQLCommentRepository {
	return &PostgreSQLCommentRepository{pool: pool}
}

func (r *PostgreSQLCommentRepository) Create(ctx context.Context, comment entities.Comment) (entities.Comment, error) {
	const query = `
		INSERT INTO comments (post_id, user_id, content)
		VALUES ($1, $2, $3)
		RETURNING id, is_visible, created_at
	`
	err := r.pool.QueryRow(ctx, query, comment.PostID, comment.UserID, comment.Content).Scan(&comment.ID, &comment.IsVisible, &comment.CreatedAt)
	if err != nil {
		return entities.Comment{}, fmt.Errorf("no se pudo crear el comentario: %w", err)
	}
	return comment, nil
}

func (r *PostgreSQLCommentRepository) FindByPostID(ctx context.Context, postID string) ([]entities.Comment, error) {
	const query = `
		SELECT id, post_id, user_id, content, toxicity_score, is_visible, created_at
		FROM comments
		WHERE post_id = $1 AND is_visible = true
		ORDER BY created_at
	`
	rows, err := r.pool.Query(ctx, query, postID)
	if err != nil {
		return nil, fmt.Errorf("no se pudieron listar los comentarios: %w", err)
	}
	defer rows.Close()

	var list []entities.Comment
	for rows.Next() {
		var c entities.Comment
		if err := rows.Scan(&c.ID, &c.PostID, &c.UserID, &c.Content, &c.ToxicityScore, &c.IsVisible, &c.CreatedAt); err != nil {
			return nil, fmt.Errorf("no se pudo leer el comentario: %w", err)
		}
		list = append(list, c)
	}
	return list, rows.Err()
}

func (r *PostgreSQLCommentRepository) Delete(ctx context.Context, id string, userID string) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM comments WHERE id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		return fmt.Errorf("no se pudo eliminar el comentario: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return repositories.ErrCommentNotFound
	}
	return nil
}
