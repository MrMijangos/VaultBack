package adapters

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"vault/src/features/assetcomments/domain/entities"
	"vault/src/features/assetcomments/domain/repositories"
)

const selectAssetCommentsQuery = `
	SELECT c.id, c.asset_id, c.user_id, c.content, c.toxicity_score, c.is_visible, c.created_at,
	       u.name, COALESCE(u.avatar_url, '')
	FROM asset_comments c
	JOIN users u ON u.id = c.user_id
`

type PostgreSQLAssetCommentRepository struct {
	pool *pgxpool.Pool
}

func NewPostgreSQLAssetCommentRepository(pool *pgxpool.Pool) *PostgreSQLAssetCommentRepository {
	return &PostgreSQLAssetCommentRepository{pool: pool}
}

func scanAssetComment(row pgx.Row) (entities.AssetComment, error) {
	var c entities.AssetComment
	err := row.Scan(
		&c.ID, &c.AssetID, &c.UserID, &c.Content, &c.ToxicityScore, &c.IsVisible, &c.CreatedAt,
		&c.AuthorName, &c.AuthorAvatarURL,
	)
	return c, err
}

func (r *PostgreSQLAssetCommentRepository) Create(ctx context.Context, comment entities.AssetComment) (entities.AssetComment, error) {
	const query = `
		INSERT INTO asset_comments (id, asset_id, user_id, content, toxicity_score, is_visible)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.pool.Exec(ctx, query,
		comment.ID, comment.AssetID, comment.UserID, comment.Content, comment.ToxicityScore, comment.IsVisible,
	)
	if err != nil {
		return entities.AssetComment{}, fmt.Errorf("no se pudo crear el comentario: %w", err)
	}
	return r.findByID(ctx, comment.ID)
}

func (r *PostgreSQLAssetCommentRepository) findByID(ctx context.Context, id string) (entities.AssetComment, error) {
	row := r.pool.QueryRow(ctx, selectAssetCommentsQuery+" WHERE c.id = $1", id)
	return scanAssetComment(row)
}

func (r *PostgreSQLAssetCommentRepository) FindByAssetID(ctx context.Context, assetID string) ([]entities.AssetComment, error) {
	rows, err := r.pool.Query(ctx, selectAssetCommentsQuery+" WHERE c.asset_id = $1 AND c.is_visible = true ORDER BY c.created_at", assetID)
	if err != nil {
		return nil, fmt.Errorf("no se pudieron listar los comentarios: %w", err)
	}
	defer rows.Close()

	var list []entities.AssetComment
	for rows.Next() {
		c, err := scanAssetComment(rows)
		if err != nil {
			return nil, fmt.Errorf("no se pudo leer el comentario: %w", err)
		}
		list = append(list, c)
	}
	return list, rows.Err()
}

func (r *PostgreSQLAssetCommentRepository) Delete(ctx context.Context, id string, userID string) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM asset_comments WHERE id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		return fmt.Errorf("no se pudo eliminar el comentario: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return repositories.ErrAssetCommentNotFound
	}
	return nil
}
