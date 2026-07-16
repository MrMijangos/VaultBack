package adapters

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"vault/src/features/comments/domain/entities"
	"vault/src/features/comments/domain/repositories"
)

// author_name/author_avatar_url no viven en comments -- se resuelven con
// un JOIN a users, igual que en posts.
const selectCommentsQuery = `
	SELECT c.id, c.post_id, c.user_id, c.content, c.toxicity_score, c.is_visible, c.created_at,
	       u.name, COALESCE(u.avatar_url, '')
	FROM comments c
	JOIN users u ON u.id = c.user_id
`

type PostgreSQLCommentRepository struct {
	pool *pgxpool.Pool
}

func NewPostgreSQLCommentRepository(pool *pgxpool.Pool) *PostgreSQLCommentRepository {
	return &PostgreSQLCommentRepository{pool: pool}
}

func scanComment(row pgx.Row) (entities.Comment, error) {
	var c entities.Comment
	err := row.Scan(
		&c.ID, &c.PostID, &c.UserID, &c.Content, &c.ToxicityScore, &c.IsVisible, &c.CreatedAt,
		&c.AuthorName, &c.AuthorAvatarURL,
	)
	return c, err
}

func (r *PostgreSQLCommentRepository) Create(ctx context.Context, comment entities.Comment) (entities.Comment, error) {
	const query = `
		INSERT INTO comments (id, post_id, user_id, content, toxicity_score, is_visible)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.pool.Exec(ctx, query,
		comment.ID, comment.PostID, comment.UserID, comment.Content, comment.ToxicityScore, comment.IsVisible,
	)
	if err != nil {
		return entities.Comment{}, fmt.Errorf("no se pudo crear el comentario: %w", err)
	}
	return r.findByID(ctx, comment.ID)
}

func (r *PostgreSQLCommentRepository) findByID(ctx context.Context, id string) (entities.Comment, error) {
	row := r.pool.QueryRow(ctx, selectCommentsQuery+" WHERE c.id = $1", id)
	return scanComment(row)
}

func (r *PostgreSQLCommentRepository) FindByPostID(ctx context.Context, postID string) ([]entities.Comment, error) {
	rows, err := r.pool.Query(ctx, selectCommentsQuery+" WHERE c.post_id = $1 AND c.is_visible = true ORDER BY c.created_at", postID)
	if err != nil {
		return nil, fmt.Errorf("no se pudieron listar los comentarios: %w", err)
	}
	defer rows.Close()

	var list []entities.Comment
	for rows.Next() {
		c, err := scanComment(rows)
		if err != nil {
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
