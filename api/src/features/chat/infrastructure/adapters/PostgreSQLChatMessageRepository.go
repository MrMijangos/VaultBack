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
	SELECT id, sender_id, recipient_id, cipher_text, encrypted_aes_key, iv, status, created_at
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
		&m.ID, &m.SenderID, &m.RecipientID, &m.CipherText, &m.EncryptedAESKey, &m.IV, &m.Status, &m.CreatedAt,
	)
	return m, err
}

func (r *PostgreSQLChatMessageRepository) Create(ctx context.Context, message entities.ChatMessage) (entities.ChatMessage, error) {
	const query = `
		INSERT INTO chat_messages (sender_id, recipient_id, cipher_text, encrypted_aes_key, iv, status)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, sender_id, recipient_id, cipher_text, encrypted_aes_key, iv, status, created_at
	`
	row := r.pool.QueryRow(ctx, query,
		message.SenderID, message.RecipientID, message.CipherText, message.EncryptedAESKey, message.IV, message.Status,
	)
	created, err := scanChatMessage(row)
	if err != nil {
		return entities.ChatMessage{}, fmt.Errorf("no se pudo enviar el mensaje: %w", err)
	}
	return created, nil
}

func (r *PostgreSQLChatMessageRepository) FindConversation(ctx context.Context, userA string, userB string) ([]entities.ChatMessage, error) {
	const query = selectChatMessagesQuery + `
		WHERE (sender_id = $1 AND recipient_id = $2) OR (sender_id = $2 AND recipient_id = $1)
		ORDER BY created_at
	`
	rows, err := r.pool.Query(ctx, query, userA, userB)
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
