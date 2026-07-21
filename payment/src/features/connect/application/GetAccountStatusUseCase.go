package application

import (
	"context"

	"vault-payment/src/core/stripeclient"
	"vault-payment/src/features/connect/domain/dto/response"
	"vault-payment/src/features/connect/domain/repositories"
)

type GetAccountStatusUseCase struct {
	accountRepo  repositories.ConnectedAccountRepository
	stripeClient stripeclient.Client
}

func NewGetAccountStatusUseCase(accountRepo repositories.ConnectedAccountRepository, stripeClient stripeclient.Client) *GetAccountStatusUseCase {
	return &GetAccountStatusUseCase{accountRepo: accountRepo, stripeClient: stripeClient}
}

// Execute devuelve nil si el vendedor nunca inició el onboarding. Si ya lo
// inició, refresca charges_enabled contra Stripe -- ese valor cambia del
// lado de Stripe conforme el vendedor va completando el KYC, no por nada
// que pase en Vault.
func (uc *GetAccountStatusUseCase) Execute(ctx context.Context, userID string) (*response.ConnectAccountResponse, error) {
	account, err := uc.accountRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if account == nil {
		return nil, nil
	}

	chargesEnabled, err := uc.stripeClient.AccountChargesEnabled(ctx, account.StripeAccountID)
	if err != nil {
		return nil, err
	}
	if chargesEnabled != account.ChargesEnabled {
		account.ChargesEnabled = chargesEnabled
		if err := uc.accountRepo.Update(ctx, account); err != nil {
			return nil, err
		}
	}

	out := response.ConnectAccountFromEntity(account)
	return &out, nil
}
