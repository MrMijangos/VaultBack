package adapters

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"vault-payment/src/features/subscriptions/domain/entities"
)

const selectSubscriptionsQuery = `
	SELECT id, user_id, plan_id, status, COALESCE(stripe_customer_id, ''), COALESCE(stripe_subscription_id, ''),
	       current_period_start, current_period_end, canceled_at, created_at
	FROM subscriptions
`

// PostgreSQLSubscriptionRepository reemplaza InMemorySubscriptionRepository.
// Un usuario puede tener varias filas a lo largo del tiempo (no hay
// UNIQUE(user_id)): CreateSubscriptionUseCase solo bloquea un Create nuevo
// si la existente sigue IsActive(), así que GetByUserID siempre trae la más
// reciente -- igual que el mapa en memoria, que sobreescribía el puntero al
// crear una nueva.
type PostgreSQLSubscriptionRepository struct {
	pool *pgxpool.Pool
}

func NewPostgreSQLSubscriptionRepository(pool *pgxpool.Pool) *PostgreSQLSubscriptionRepository {
	return &PostgreSQLSubscriptionRepository{pool: pool}
}

func scanSubscription(row pgx.Row) (*entities.Subscription, error) {
	var s entities.Subscription
	err := row.Scan(
		&s.ID, &s.UserID, &s.PlanID, &s.Status, &s.StripeCustomerID, &s.StripeSubscriptionID,
		&s.CurrentPeriodStart, &s.CurrentPeriodEnd, &s.CanceledAt, &s.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *PostgreSQLSubscriptionRepository) Create(ctx context.Context, sub *entities.Subscription) error {
	const query = `
		INSERT INTO subscriptions (id, user_id, plan_id, status, stripe_customer_id, stripe_subscription_id, current_period_start, current_period_end, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.pool.Exec(ctx, query,
		sub.ID, sub.UserID, sub.PlanID, sub.Status, sub.StripeCustomerID, sub.StripeSubscriptionID,
		sub.CurrentPeriodStart, sub.CurrentPeriodEnd, sub.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("no se pudo crear la suscripción: %w", err)
	}
	return nil
}

func (r *PostgreSQLSubscriptionRepository) Update(ctx context.Context, sub *entities.Subscription) error {
	const query = `
		UPDATE subscriptions
		SET status = $1, stripe_customer_id = $2, stripe_subscription_id = $3,
		    current_period_start = $4, current_period_end = $5, canceled_at = $6
		WHERE id = $7
	`
	tag, err := r.pool.Exec(ctx, query,
		sub.Status, sub.StripeCustomerID, sub.StripeSubscriptionID,
		sub.CurrentPeriodStart, sub.CurrentPeriodEnd, sub.CanceledAt,
		sub.ID,
	)
	if err != nil {
		return fmt.Errorf("no se pudo actualizar la suscripción: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("suscripción %q no existe", sub.ID)
	}
	return nil
}

func (r *PostgreSQLSubscriptionRepository) GetByUserID(ctx context.Context, userID string) (*entities.Subscription, error) {
	row := r.pool.QueryRow(ctx, selectSubscriptionsQuery+" WHERE user_id = $1 ORDER BY created_at DESC LIMIT 1", userID)
	sub, err := scanSubscription(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("no se pudo obtener la suscripción: %w", err)
	}
	return sub, nil
}

func (r *PostgreSQLSubscriptionRepository) GetByStripeSubscriptionID(ctx context.Context, stripeSubscriptionID string) (*entities.Subscription, error) {
	row := r.pool.QueryRow(ctx, selectSubscriptionsQuery+" WHERE stripe_subscription_id = $1", stripeSubscriptionID)
	sub, err := scanSubscription(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("no se pudo obtener la suscripción: %w", err)
	}
	return sub, nil
}
