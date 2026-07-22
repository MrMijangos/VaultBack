package adapters

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"vault-payment/src/features/connect/domain/entities"
)

const selectConnectedAccountsQuery = `
	SELECT user_id, stripe_account_id, charges_enabled, created_at
	FROM connected_accounts
`

// PostgreSQLConnectedAccountRepository reemplaza
// InMemoryConnectedAccountRepository. La entidad no tiene ID propio -- es
// 1:1 con el usuario (misma clave primaria), así que Create hace upsert para
// igualar el "siempre sobreescribe" del mapa en memoria.
type PostgreSQLConnectedAccountRepository struct {
	pool *pgxpool.Pool
}

func NewPostgreSQLConnectedAccountRepository(pool *pgxpool.Pool) *PostgreSQLConnectedAccountRepository {
	return &PostgreSQLConnectedAccountRepository{pool: pool}
}

func scanConnectedAccount(row pgx.Row) (*entities.ConnectedAccount, error) {
	var a entities.ConnectedAccount
	err := row.Scan(&a.UserID, &a.StripeAccountID, &a.ChargesEnabled, &a.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *PostgreSQLConnectedAccountRepository) Create(ctx context.Context, account *entities.ConnectedAccount) error {
	const query = `
		INSERT INTO connected_accounts (user_id, stripe_account_id, charges_enabled, created_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id) DO UPDATE
		SET stripe_account_id = EXCLUDED.stripe_account_id, charges_enabled = EXCLUDED.charges_enabled
	`
	_, err := r.pool.Exec(ctx, query, account.UserID, account.StripeAccountID, account.ChargesEnabled, account.CreatedAt)
	if err != nil {
		return fmt.Errorf("no se pudo crear la cuenta conectada: %w", err)
	}
	return nil
}

func (r *PostgreSQLConnectedAccountRepository) Update(ctx context.Context, account *entities.ConnectedAccount) error {
	const query = `
		UPDATE connected_accounts
		SET stripe_account_id = $1, charges_enabled = $2
		WHERE user_id = $3
	`
	tag, err := r.pool.Exec(ctx, query, account.StripeAccountID, account.ChargesEnabled, account.UserID)
	if err != nil {
		return fmt.Errorf("no se pudo actualizar la cuenta conectada: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("el usuario %q no tiene una cuenta conectada", account.UserID)
	}
	return nil
}

func (r *PostgreSQLConnectedAccountRepository) GetByUserID(ctx context.Context, userID string) (*entities.ConnectedAccount, error) {
	row := r.pool.QueryRow(ctx, selectConnectedAccountsQuery+" WHERE user_id = $1", userID)
	account, err := scanConnectedAccount(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("no se pudo obtener la cuenta conectada: %w", err)
	}
	return account, nil
}
