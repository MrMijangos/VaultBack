package adapters

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgreSQLFCMTokenRepository struct {
	pool *pgxpool.Pool
}

func NewPostgreSQLFCMTokenRepository(pool *pgxpool.Pool) *PostgreSQLFCMTokenRepository {
	return &PostgreSQLFCMTokenRepository{pool: pool}
}

// Register hace upsert sobre el UNIQUE de token: si es nuevo lo inserta: si
// ya existía -- del mismo usuario (llamada repetida, p.ej. onTokenRefresh) o
// de otro (reinstalación con otra cuenta en el mismo dispositivo) -- lo
// reasigna a userID en vez de fallar o dejarlo huérfano.
func (r *PostgreSQLFCMTokenRepository) Register(ctx context.Context, userID string, token string, platform *string) error {
	const query = `
		INSERT INTO fcm_tokens (user_id, token, platform)
		VALUES ($1, $2, $3)
		ON CONFLICT (token) DO UPDATE
		SET user_id = EXCLUDED.user_id, platform = EXCLUDED.platform, updated_at = now()
	`
	if _, err := r.pool.Exec(ctx, query, userID, token, platform); err != nil {
		return fmt.Errorf("no se pudo registrar el token FCM: %w", err)
	}
	return nil
}

func (r *PostgreSQLFCMTokenRepository) FindTokensByUserID(ctx context.Context, userID string) ([]string, error) {
	rows, err := r.pool.Query(ctx, `SELECT token FROM fcm_tokens WHERE user_id = $1`, userID)
	if err != nil {
		return nil, fmt.Errorf("no se pudieron obtener los tokens FCM: %w", err)
	}
	defer rows.Close()

	var tokens []string
	for rows.Next() {
		var token string
		if err := rows.Scan(&token); err != nil {
			return nil, fmt.Errorf("no se pudo leer un token FCM: %w", err)
		}
		tokens = append(tokens, token)
	}
	return tokens, rows.Err()
}

func (r *PostgreSQLFCMTokenRepository) DeleteToken(ctx context.Context, token string) error {
	if _, err := r.pool.Exec(ctx, `DELETE FROM fcm_tokens WHERE token = $1`, token); err != nil {
		return fmt.Errorf("no se pudo eliminar el token FCM: %w", err)
	}
	return nil
}
