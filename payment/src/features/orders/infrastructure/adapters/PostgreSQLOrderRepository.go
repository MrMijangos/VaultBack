package adapters

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"vault-payment/src/features/orders/domain/entities"
)

const selectOrdersQuery = `
	SELECT id, buyer_id, seller_id, asset_id, amount_cents, commission_cents, seller_amount_cents,
	       currency, status, COALESCE(stripe_customer_id, ''), COALESCE(stripe_payment_intent_id, ''),
	       COALESCE(stripe_transfer_id, ''), created_at, confirmed_at
	FROM orders
`

// PostgreSQLOrderRepository reemplaza InMemoryOrderRepository.
type PostgreSQLOrderRepository struct {
	pool *pgxpool.Pool
}

func NewPostgreSQLOrderRepository(pool *pgxpool.Pool) *PostgreSQLOrderRepository {
	return &PostgreSQLOrderRepository{pool: pool}
}

func scanOrder(row pgx.Row) (*entities.Order, error) {
	var o entities.Order
	err := row.Scan(
		&o.ID, &o.BuyerID, &o.SellerID, &o.AssetID, &o.AmountCents, &o.CommissionCents, &o.SellerAmountCents,
		&o.Currency, &o.Status, &o.StripeCustomerID, &o.StripePaymentIntentID,
		&o.StripeTransferID, &o.CreatedAt, &o.ConfirmedAt,
	)
	if err != nil {
		return nil, err
	}
	return &o, nil
}

func (r *PostgreSQLOrderRepository) Create(ctx context.Context, order *entities.Order) error {
	const query = `
		INSERT INTO orders (id, buyer_id, seller_id, asset_id, amount_cents, commission_cents, seller_amount_cents, currency, status, stripe_customer_id, stripe_payment_intent_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`
	_, err := r.pool.Exec(ctx, query,
		order.ID, order.BuyerID, order.SellerID, order.AssetID, order.AmountCents, order.CommissionCents,
		order.SellerAmountCents, order.Currency, order.Status, order.StripeCustomerID, order.StripePaymentIntentID,
		order.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("no se pudo crear la orden: %w", err)
	}
	return nil
}

func (r *PostgreSQLOrderRepository) Update(ctx context.Context, order *entities.Order) error {
	const query = `
		UPDATE orders
		SET status = $1, stripe_transfer_id = $2, confirmed_at = $3
		WHERE id = $4
	`
	tag, err := r.pool.Exec(ctx, query, order.Status, order.StripeTransferID, order.ConfirmedAt, order.ID)
	if err != nil {
		return fmt.Errorf("no se pudo actualizar la orden: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("orden %q no existe", order.ID)
	}
	return nil
}

func (r *PostgreSQLOrderRepository) GetByID(ctx context.Context, id string) (*entities.Order, error) {
	row := r.pool.QueryRow(ctx, selectOrdersQuery+" WHERE id = $1", id)
	order, err := scanOrder(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("no se pudo obtener la orden: %w", err)
	}
	return order, nil
}

func (r *PostgreSQLOrderRepository) ListByBuyerID(ctx context.Context, buyerID string) ([]*entities.Order, error) {
	rows, err := r.pool.Query(ctx, selectOrdersQuery+" WHERE buyer_id = $1 ORDER BY created_at DESC", buyerID)
	if err != nil {
		return nil, fmt.Errorf("no se pudieron listar las ordenes: %w", err)
	}
	defer rows.Close()

	var out []*entities.Order
	for rows.Next() {
		o, err := scanOrder(rows)
		if err != nil {
			return nil, fmt.Errorf("no se pudo leer la orden: %w", err)
		}
		out = append(out, o)
	}
	return out, rows.Err()
}

func (r *PostgreSQLOrderRepository) ListBySellerID(ctx context.Context, sellerID string) ([]*entities.Order, error) {
	rows, err := r.pool.Query(ctx, selectOrdersQuery+" WHERE seller_id = $1 ORDER BY created_at DESC", sellerID)
	if err != nil {
		return nil, fmt.Errorf("no se pudieron listar las ordenes: %w", err)
	}
	defer rows.Close()

	var out []*entities.Order
	for rows.Next() {
		o, err := scanOrder(rows)
		if err != nil {
			return nil, fmt.Errorf("no se pudo leer la orden: %w", err)
		}
		out = append(out, o)
	}
	return out, rows.Err()
}
