package application

import (
	"context"
	"fmt"
	"time"

	"vault-payment/src/core/stripeclient"
	"vault-payment/src/features/connect/domain/entities"
	"vault-payment/src/features/connect/domain/repositories"
)

// defaultAccountCountry -- Vault opera en México (los planes están en MXN);
// se deja fija por ahora, no hay selección de país en el spec.
const defaultAccountCountry = "MX"

type CreateOnboardingLinkUseCase struct {
	accountRepo  repositories.ConnectedAccountRepository
	stripeClient stripeclient.Client
}

func NewCreateOnboardingLinkUseCase(accountRepo repositories.ConnectedAccountRepository, stripeClient stripeclient.Client) *CreateOnboardingLinkUseCase {
	return &CreateOnboardingLinkUseCase{accountRepo: accountRepo, stripeClient: stripeClient}
}

// Execute crea la cuenta de Stripe Connect del vendedor si es su primera
// vez, y siempre devuelve un link de onboarding nuevo -- los AccountLink de
// Stripe son de un solo uso y expiran rápido, así que regenerarlo en cada
// llamada es el comportamiento esperado, no un error.
func (uc *CreateOnboardingLinkUseCase) Execute(ctx context.Context, userID, email, refreshURL, returnURL string) (string, error) {
	account, err := uc.accountRepo.GetByUserID(ctx, userID)
	if err != nil {
		return "", err
	}

	if account == nil {
		accountID, err := uc.stripeClient.CreateExpressAccount(ctx, email, defaultAccountCountry)
		if err != nil {
			return "", fmt.Errorf("no se pudo crear la cuenta de Stripe Connect: %w", err)
		}

		account = &entities.ConnectedAccount{
			UserID:          userID,
			StripeAccountID: accountID,
			ChargesEnabled:  false,
			CreatedAt:       time.Now().UTC(),
		}
		if err := uc.accountRepo.Create(ctx, account); err != nil {
			return "", err
		}
	}

	link, err := uc.stripeClient.CreateOnboardingLink(ctx, account.StripeAccountID, refreshURL, returnURL)
	if err != nil {
		return "", fmt.Errorf("no se pudo generar el link de onboarding: %w", err)
	}
	return link, nil
}
