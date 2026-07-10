package adapters

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"vault/src/features/notifications/domain/entities"
	"vault/src/features/notifications/domain/repositories"
)

const selectNotificationsQuery = `
	SELECT id, user_id, type, subtype, title, body, data, read, created_at
	FROM notifications
`

type PostgreSQLNotificationRepository struct {
	pool *pgxpool.Pool
}

func NewPostgreSQLNotificationRepository(pool *pgxpool.Pool) *PostgreSQLNotificationRepository {
	return &PostgreSQLNotificationRepository{pool: pool}
}

func scanNotification(row pgx.Row) (entities.Notification, error) {
	var n entities.Notification
	err := row.Scan(&n.ID, &n.UserID, &n.Type, &n.Subtype, &n.Title, &n.Body, &n.Data, &n.Read, &n.CreatedAt)
	return n, err
}

func (r *PostgreSQLNotificationRepository) Create(ctx context.Context, notification entities.Notification) (entities.Notification, error) {
	const query = `
		INSERT INTO notifications (user_id, type, subtype, title, body, data)
		VALUES ($1, $2, $3, $4, $5, COALESCE($6, '{}'::jsonb))
		RETURNING id, data, read, created_at
	`
	err := r.pool.QueryRow(ctx, query,
		notification.UserID, notification.Type, notification.Subtype, notification.Title, notification.Body, notification.Data,
	).Scan(&notification.ID, &notification.Data, &notification.Read, &notification.CreatedAt)
	if err != nil {
		return entities.Notification{}, fmt.Errorf("no se pudo crear la notificacion: %w", err)
	}
	return notification, nil
}

func (r *PostgreSQLNotificationRepository) FindByUserID(ctx context.Context, userID string) ([]entities.Notification, error) {
	rows, err := r.pool.Query(ctx, selectNotificationsQuery+" WHERE user_id = $1 ORDER BY created_at DESC", userID)
	if err != nil {
		return nil, fmt.Errorf("no se pudieron listar las notificaciones: %w", err)
	}
	defer rows.Close()

	var list []entities.Notification
	for rows.Next() {
		n, err := scanNotification(rows)
		if err != nil {
			return nil, fmt.Errorf("no se pudo leer la notificacion: %w", err)
		}
		list = append(list, n)
	}
	return list, rows.Err()
}

func (r *PostgreSQLNotificationRepository) MarkAsRead(ctx context.Context, id string, userID string) (entities.Notification, error) {
	tag, err := r.pool.Exec(ctx, `UPDATE notifications SET read = true WHERE id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		return entities.Notification{}, fmt.Errorf("no se pudo marcar como leida: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return entities.Notification{}, repositories.ErrNotificationNotFound
	}

	row := r.pool.QueryRow(ctx, selectNotificationsQuery+" WHERE id = $1", id)
	n, err := scanNotification(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return entities.Notification{}, repositories.ErrNotificationNotFound
	}
	if err != nil {
		return entities.Notification{}, fmt.Errorf("no se pudo obtener la notificacion: %w", err)
	}
	return n, nil
}

func (r *PostgreSQLNotificationRepository) Delete(ctx context.Context, id string, userID string) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM notifications WHERE id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		return fmt.Errorf("no se pudo eliminar la notificacion: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return repositories.ErrNotificationNotFound
	}
	return nil
}
