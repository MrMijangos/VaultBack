package adapters

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"vault/src/features/chat/domain/entities"
	"vault/src/features/chat/domain/repositories"
)

const selectChatMessagesQuery = `
	SELECT id, sender_id, recipient_id, cipher_text, encrypted_aes_key,
	       COALESCE(encrypted_aes_key_sender, ''), iv, status, created_at
	FROM chat_messages
`

type PostgreSQLChatMessageRepository struct {
	pool *pgxpool.Pool
}

func NewPostgreSQLChatMessageRepository(pool *pgxpool.Pool) *PostgreSQLChatMessageRepository {
	return &PostgreSQLChatMessageRepository{pool: pool}
}

func scanChatMessage(row pgx.Row) (entities.ChatMessage, error) {
	var m entities.ChatMessage
	err := row.Scan(
		&m.ID, &m.SenderID, &m.RecipientID, &m.CipherText, &m.EncryptedAESKey,
		&m.EncryptedAESKeySender, &m.IV, &m.Status, &m.CreatedAt,
	)
	return m, err
}

func (r *PostgreSQLChatMessageRepository) Create(ctx context.Context, message entities.ChatMessage) (entities.ChatMessage, error) {
	const query = `
		INSERT INTO chat_messages (sender_id, recipient_id, cipher_text, encrypted_aes_key, encrypted_aes_key_sender, iv, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, sender_id, recipient_id, cipher_text, encrypted_aes_key,
		          COALESCE(encrypted_aes_key_sender, ''), iv, status, created_at
	`
	row := r.pool.QueryRow(ctx, query,
		message.SenderID, message.RecipientID, message.CipherText, message.EncryptedAESKey,
		message.EncryptedAESKeySender, message.IV, message.Status,
	)
	created, err := scanChatMessage(row)
	if err != nil {
		return entities.ChatMessage{}, fmt.Errorf("no se pudo enviar el mensaje: %w", err)
	}
	return created, nil
}

func (r *PostgreSQLChatMessageRepository) FindConversation(ctx context.Context, meID string, otherID string) ([]entities.ChatMessage, error) {
	const query = selectChatMessagesQuery + `
		WHERE ((sender_id = $1 AND recipient_id = $2) OR (sender_id = $2 AND recipient_id = $1))
		  AND NOT ((sender_id = $1 AND deleted_by_sender) OR (recipient_id = $1 AND deleted_by_recipient))
		ORDER BY created_at
	`
	rows, err := r.pool.Query(ctx, query, meID, otherID)
	if err != nil {
		return nil, fmt.Errorf("no se pudo obtener la conversacion: %w", err)
	}
	defer rows.Close()

	var list []entities.ChatMessage
	for rows.Next() {
		m, err := scanChatMessage(rows)
		if err != nil {
			return nil, fmt.Errorf("no se pudo leer el mensaje: %w", err)
		}
		list = append(list, m)
	}
	return list, rows.Err()
}

func (r *PostgreSQLChatMessageRepository) FindConversationsForUser(ctx context.Context, userID string) ([]entities.ConversationSummary, error) {
	// conv: el último mensaje con cada contraparte (ROW_NUMBER = 1).
	// unread: cuántos mensajes de cada contraparte siguen sin leerse.
	const query = `
		WITH conv AS (
			SELECT
				CASE WHEN sender_id = $1 THEN recipient_id ELSE sender_id END AS other_user_id,
				id, sender_id, recipient_id, cipher_text, encrypted_aes_key,
				COALESCE(encrypted_aes_key_sender, '') AS encrypted_aes_key_sender,
				iv, status, created_at,
				ROW_NUMBER() OVER (
					PARTITION BY CASE WHEN sender_id = $1 THEN recipient_id ELSE sender_id END
					ORDER BY created_at DESC
				) AS rn
			FROM chat_messages
			WHERE (sender_id = $1 OR recipient_id = $1)
			  AND NOT ((sender_id = $1 AND deleted_by_sender) OR (recipient_id = $1 AND deleted_by_recipient))
		),
		unread AS (
			SELECT sender_id AS other_user_id, COUNT(*) AS unread_count
			FROM chat_messages
			WHERE recipient_id = $1 AND status <> 'read' AND NOT deleted_by_recipient
			GROUP BY sender_id
		)
		SELECT
			conv.other_user_id, u.name, COALESCE(u.avatar_url, ''),
			conv.id, conv.sender_id, conv.recipient_id, conv.cipher_text, conv.encrypted_aes_key,
			conv.encrypted_aes_key_sender, conv.iv, conv.status, conv.created_at,
			COALESCE(unread.unread_count, 0)
		FROM conv
		JOIN users u ON u.id = conv.other_user_id
		LEFT JOIN unread ON unread.other_user_id = conv.other_user_id
		WHERE conv.rn = 1
		ORDER BY conv.created_at DESC
	`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("no se pudieron listar las conversaciones: %w", err)
	}
	defer rows.Close()

	var list []entities.ConversationSummary
	for rows.Next() {
		var s entities.ConversationSummary
		if err := rows.Scan(
			&s.OtherUserID, &s.OtherUserName, &s.OtherUserAvatarURL,
			&s.LastMessage.ID, &s.LastMessage.SenderID, &s.LastMessage.RecipientID,
			&s.LastMessage.CipherText, &s.LastMessage.EncryptedAESKey,
			&s.LastMessage.EncryptedAESKeySender, &s.LastMessage.IV,
			&s.LastMessage.Status, &s.LastMessage.CreatedAt,
			&s.UnreadCount,
		); err != nil {
			return nil, fmt.Errorf("no se pudo leer la conversacion: %w", err)
		}
		list = append(list, s)
	}
	return list, rows.Err()
}

func (r *PostgreSQLChatMessageRepository) UpdateStatus(ctx context.Context, id string, recipientID string, status string) (entities.ChatMessage, error) {
	const query = `UPDATE chat_messages SET status = $1 WHERE id = $2 AND recipient_id = $3`

	tag, err := r.pool.Exec(ctx, query, status, id, recipientID)
	if err != nil {
		return entities.ChatMessage{}, fmt.Errorf("no se pudo actualizar el estado del mensaje: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return entities.ChatMessage{}, repositories.ErrChatMessageNotFound
	}

	row := r.pool.QueryRow(ctx, selectChatMessagesQuery+" WHERE id = $1", id)
	updated, err := scanChatMessage(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return entities.ChatMessage{}, repositories.ErrChatMessageNotFound
	}
	if err != nil {
		return entities.ChatMessage{}, fmt.Errorf("no se pudo leer el mensaje actualizado: %w", err)
	}
	return updated, nil
}

func (r *PostgreSQLChatMessageRepository) DeleteMessage(ctx context.Context, id string, userID string) error {
	const query = `
		UPDATE chat_messages
		SET deleted_by_sender = CASE WHEN sender_id = $2 THEN true ELSE deleted_by_sender END,
		    deleted_by_recipient = CASE WHEN recipient_id = $2 THEN true ELSE deleted_by_recipient END
		WHERE id = $1 AND (sender_id = $2 OR recipient_id = $2)
	`
	tag, err := r.pool.Exec(ctx, query, id, userID)
	if err != nil {
		return fmt.Errorf("no se pudo eliminar el mensaje: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return repositories.ErrChatMessageNotFound
	}
	return nil
}

func (r *PostgreSQLChatMessageRepository) DeleteConversation(ctx context.Context, userID string, otherID string) error {
	const query = `
		UPDATE chat_messages
		SET deleted_by_sender = CASE WHEN sender_id = $1 THEN true ELSE deleted_by_sender END,
		    deleted_by_recipient = CASE WHEN recipient_id = $1 THEN true ELSE deleted_by_recipient END
		WHERE (sender_id = $1 AND recipient_id = $2) OR (sender_id = $2 AND recipient_id = $1)
	`
	if _, err := r.pool.Exec(ctx, query, userID, otherID); err != nil {
		return fmt.Errorf("no se pudo eliminar la conversacion: %w", err)
	}
	return nil
}
