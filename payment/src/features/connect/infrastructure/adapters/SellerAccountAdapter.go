package adapters

import (
	"context"

	"vault-payment/src/features/connect/domain/repositories"
)

// SellerAccountAdapter implementa orders/domain/repositories.SellerAccountProvider
// (tipado estructural, sin import cruzado entre features).
type SellerAccountAdapter struct {
	accountRepo repositories.ConnectedAccountRepository
}

func NewSellerAccountAdapter(accountRepo repositories.ConnectedAccountRepository) *SellerAccountAdapter {
	return &SellerAccountAdapter{accountRepo: accountRepo}
}

func (a *SellerAccountAdapter) GetChargesEnabledAccountID(ctx context.Context, sellerID string) (string, bool, error) {
	account, err := a.accountRepo.GetByUserID(ctx, sellerID)
	if err != nil {
		return "", false, err
	}
	if account == nil {
		return "", false, nil
	}
	return account.StripeAccountID, account.ChargesEnabled, nil
}
