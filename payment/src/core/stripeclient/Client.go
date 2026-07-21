package stripeclient

import (
	"context"
	"errors"
	"time"

	"github.com/stripe/stripe-go/v86"
)

var ErrNotConfigured = errors.New("Stripe todavia no esta configurado en este servicio (falta STRIPE_SECRET_KEY)")

// IsNotConfigured ayuda a los controllers a distinguir "falta STRIPE_SECRET_KEY"
// (503, se resuelve configurando la cuenta) de cualquier otro error real de
// Stripe (500). Funciona incluso si el error llegó envuelto con %w.
func IsNotConfigured(err error) bool {
	return errors.Is(err, ErrNotConfigured)
}

// Client es la interfaz que usan los casos de uso -- así CreateSubscriptionUseCase
// y CancelSubscriptionUseCase no dependen directamente del SDK de Stripe.
type Client interface {
	Configured() bool
	// CreateCustomerWithPaymentMethod crea (o reutiliza) el customer de Stripe,
	// adjunta el método de pago ya tokenizado del lado del cliente (flutter_stripe
	// nunca manda el número de tarjeta a este backend, solo el payment_method_id)
	// y lo deja como default para cobros automáticos. Devuelve el ID del método
	// de pago tal como quedó adjuntado al customer -- con los payment methods de
	// prueba de Stripe (ej. "pm_card_visa") el ID adjuntado real es distinto al
	// que se mandó, así que las llamadas siguientes (ChargeBuyer, etc.) deben
	// usar el valor devuelto aquí, no el que llegó del cliente.
	CreateCustomerWithPaymentMethod(ctx context.Context, email string, paymentMethodID string) (customerID string, attachedPaymentMethodID string, err error)
	// CreateSubscription arma la suscripción con un solo price (un plan = un
	// price de Stripe) y devuelve el periodo de facturación vigente.
	CreateSubscription(ctx context.Context, customerID string, priceID string) (subscriptionID string, periodStart time.Time, periodEnd time.Time, err error)
	CancelSubscription(ctx context.Context, subscriptionID string) error
	ConstructWebhookEvent(payload []byte, sigHeader string, webhookSecret string) (stripe.Event, error)

	// -- Escrow de compra-venta (orders/) y onboarding de vendedores (connect/) --

	// CreateExpressAccount crea la cuenta de Stripe Connect Express del
	// vendedor. El KYC (identidad, datos bancarios) lo recolecta Stripe en
	// su propio flujo hospedado -- ver CreateOnboardingLink -- Vault nunca
	// ve esos datos.
	CreateExpressAccount(ctx context.Context, email string, country string) (accountID string, err error)
	// CreateOnboardingLink genera la URL de un solo uso a la que se
	// redirige al vendedor para completar el onboarding de Stripe Connect.
	CreateOnboardingLink(ctx context.Context, accountID string, refreshURL string, returnURL string) (url string, err error)
	// AccountChargesEnabled refleja si Stripe ya validó al vendedor lo
	// suficiente como para poder recibir transferencias (Transfer) --
	// mientras sea false, no se le puede transferir el pago liberado.
	AccountChargesEnabled(ctx context.Context, accountID string) (bool, error)
	// ChargeBuyer cobra el monto completo de una orden al comprador y lo
	// deja en el balance de la cuenta de Stripe de Vault (no en la del
	// vendedor) -- así es como se "retiene" el dinero: el patrón de Stripe
	// Connect "separate charges and transfers".
	ChargeBuyer(ctx context.Context, customerID string, paymentMethodID string, amountCents int64, currency string) (paymentIntentID string, err error)
	// ReleaseToSeller transfiere el monto (ya con la comisión de Vault
	// descontada) desde el balance de Vault hacia la cuenta conectada del
	// vendedor. Se llama solo cuando el comprador confirma haber recibido
	// el producto.
	ReleaseToSeller(ctx context.Context, destinationAccountID string, amountCents int64, currency string, transferGroup string) (transferID string, err error)
}

type stripeClient struct {
	sc *stripe.Client
}

func New(secretKey string) Client {
	if secretKey == "" {
		return notConfiguredClient{}
	}
	return &stripeClient{sc: stripe.NewClient(secretKey)}
}

func (c *stripeClient) Configured() bool { return true }

func (c *stripeClient) CreateCustomerWithPaymentMethod(ctx context.Context, email string, paymentMethodID string) (string, string, error) {
	customer, err := c.sc.V1Customers.Create(ctx, &stripe.CustomerCreateParams{
		Email: strPtr(email),
	})
	if err != nil {
		return "", "", err
	}

	attached, err := c.sc.V1PaymentMethods.Attach(ctx, paymentMethodID, &stripe.PaymentMethodAttachParams{
		Customer: strPtr(customer.ID),
	})
	if err != nil {
		return "", "", err
	}

	if _, err := c.sc.V1Customers.Update(ctx, customer.ID, &stripe.CustomerUpdateParams{
		InvoiceSettings: &stripe.CustomerUpdateInvoiceSettingsParams{
			DefaultPaymentMethod: strPtr(attached.ID),
		},
	}); err != nil {
		return "", "", err
	}

	return customer.ID, attached.ID, nil
}

func (c *stripeClient) CreateSubscription(ctx context.Context, customerID string, priceID string) (string, time.Time, time.Time, error) {
	// "default_incomplete" deja el primer cobro sin confirmar a propósito,
	// esperando que un cliente con interfaz (Stripe.js/flutter_stripe)
	// confirme el PaymentIntent -- como este backend no tiene ese paso,
	// usamos "error_if_incomplete" + OffSession para que Stripe intente el
	// cobro inmediatamente con el método de pago que ya dejamos como
	// default (ver CreateCustomerWithPaymentMethod). Si el cobro falla, la
	// llamada regresa error y no se crea la suscripción -- si tiene éxito,
	// la suscripción queda "active" de una vez, sin quedar "incomplete".
	sub, err := c.sc.V1Subscriptions.Create(ctx, &stripe.SubscriptionCreateParams{
		Customer: strPtr(customerID),
		Items: []*stripe.SubscriptionCreateItemParams{
			{Price: strPtr(priceID)},
		},
		PaymentBehavior: strPtr("error_if_incomplete"),
		OffSession:      stripe.Bool(true),
	})
	if err != nil {
		return "", time.Time{}, time.Time{}, err
	}

	var periodStart, periodEnd time.Time
	if sub.Items != nil && len(sub.Items.Data) > 0 {
		item := sub.Items.Data[0]
		periodStart = time.Unix(item.CurrentPeriodStart, 0).UTC()
		periodEnd = time.Unix(item.CurrentPeriodEnd, 0).UTC()
	}

	return sub.ID, periodStart, periodEnd, nil
}

func (c *stripeClient) CancelSubscription(ctx context.Context, subscriptionID string) error {
	_, err := c.sc.V1Subscriptions.Cancel(ctx, subscriptionID, &stripe.SubscriptionCancelParams{})
	return err
}

func (c *stripeClient) ConstructWebhookEvent(payload []byte, sigHeader string, webhookSecret string) (stripe.Event, error) {
	return stripe.ConstructEvent(payload, sigHeader, webhookSecret)
}

// El vendedor se modela como cuenta "Recipient" (no "Merchant"): no es
// merchant of record, solo recibe fondos vía Transfer -- así lo pide el
// patrón "separate charges and transfers". Vault (la plataforma) se queda
// con las comisiones y la responsabilidad de pérdidas (fees_collector /
// losses_collector = "application"), como corresponde a un marketplace.
func (c *stripeClient) CreateExpressAccount(ctx context.Context, email string, country string) (string, error) {
	account, err := c.sc.V2CoreAccounts.Create(ctx, &stripe.V2CoreAccountCreateParams{
		ContactEmail: strPtr(email),
		Dashboard:    strPtr("express"),
		Identity: &stripe.V2CoreAccountCreateIdentityParams{
			Country: strPtr(country),
		},
		Configuration: &stripe.V2CoreAccountCreateConfigurationParams{
			Recipient: &stripe.V2CoreAccountCreateConfigurationRecipientParams{
				Capabilities: &stripe.V2CoreAccountCreateConfigurationRecipientCapabilitiesParams{
					StripeBalance: &stripe.V2CoreAccountCreateConfigurationRecipientCapabilitiesStripeBalanceParams{
						StripeTransfers: &stripe.V2CoreAccountCreateConfigurationRecipientCapabilitiesStripeBalanceStripeTransfersParams{
							Requested: stripe.Bool(true),
						},
					},
				},
			},
		},
		Defaults: &stripe.V2CoreAccountCreateDefaultsParams{
			Responsibilities: &stripe.V2CoreAccountCreateDefaultsResponsibilitiesParams{
				FeesCollector:   strPtr("application"),
				LossesCollector: strPtr("application"),
			},
		},
	})
	if err != nil {
		return "", err
	}
	return account.ID, nil
}

func (c *stripeClient) CreateOnboardingLink(ctx context.Context, accountID string, refreshURL string, returnURL string) (string, error) {
	link, err := c.sc.V2CoreAccountLinks.Create(ctx, &stripe.V2CoreAccountLinkCreateParams{
		Account: strPtr(accountID),
		UseCase: &stripe.V2CoreAccountLinkCreateUseCaseParams{
			Type: strPtr("account_onboarding"),
			AccountOnboarding: &stripe.V2CoreAccountLinkCreateUseCaseAccountOnboardingParams{
				Configurations: stripe.StringSlice([]string{"recipient"}),
				RefreshURL:     strPtr(refreshURL),
				ReturnURL:      strPtr(returnURL),
			},
		},
	})
	if err != nil {
		return "", err
	}
	return link.URL, nil
}

// AccountChargesEnabled ya no lee el campo `charges_enabled` (deprecado en
// v1) sino el estado de la capability stripe_transfers de la configuración
// Recipient -- es el reemplazo correcto en la API v2 para saber si el
// vendedor ya puede recibir un Transfer.
func (c *stripeClient) AccountChargesEnabled(ctx context.Context, accountID string) (bool, error) {
	account, err := c.sc.V2CoreAccounts.Retrieve(ctx, accountID, &stripe.V2CoreAccountRetrieveParams{
		Include: stripe.StringSlice([]string{"configuration.recipient"}),
	})
	if err != nil {
		return false, err
	}
	if account.Configuration == nil || account.Configuration.Recipient == nil ||
		account.Configuration.Recipient.Capabilities == nil ||
		account.Configuration.Recipient.Capabilities.StripeBalance == nil ||
		account.Configuration.Recipient.Capabilities.StripeBalance.StripeTransfers == nil {
		return false, nil
	}
	status := account.Configuration.Recipient.Capabilities.StripeBalance.StripeTransfers.Status
	return status == stripe.V2CoreAccountConfigurationRecipientCapabilitiesStripeBalanceStripeTransfersStatusActive, nil
}

func (c *stripeClient) ChargeBuyer(ctx context.Context, customerID string, paymentMethodID string, amountCents int64, currency string) (string, error) {
	pi, err := c.sc.V1PaymentIntents.Create(ctx, &stripe.PaymentIntentCreateParams{
		Amount:        stripe.Int64(amountCents),
		Currency:      strPtr(currency),
		Customer:      strPtr(customerID),
		PaymentMethod: strPtr(paymentMethodID),
		Confirm:       stripe.Bool(true),
		OffSession:    stripe.Bool(true),
	})
	if err != nil {
		return "", err
	}
	return pi.ID, nil
}

func (c *stripeClient) ReleaseToSeller(ctx context.Context, destinationAccountID string, amountCents int64, currency string, transferGroup string) (string, error) {
	transfer, err := c.sc.V1Transfers.Create(ctx, &stripe.TransferCreateParams{
		Amount:        stripe.Int64(amountCents),
		Currency:      strPtr(currency),
		Destination:   strPtr(destinationAccountID),
		TransferGroup: strPtr(transferGroup),
	})
	if err != nil {
		return "", err
	}
	return transfer.ID, nil
}

// notConfiguredClient se usa mientras no exista STRIPE_SECRET_KEY (todavía no
// hay cuenta de Stripe) -- el servicio arranca igual, y las rutas que sí
// necesitan Stripe devuelven ErrNotConfigured en vez de un panic o un 500
// críptico.
type notConfiguredClient struct{}

func (notConfiguredClient) Configured() bool { return false }

func (notConfiguredClient) CreateCustomerWithPaymentMethod(context.Context, string, string) (string, string, error) {
	return "", "", ErrNotConfigured
}

func (notConfiguredClient) CreateSubscription(context.Context, string, string) (string, time.Time, time.Time, error) {
	return "", time.Time{}, time.Time{}, ErrNotConfigured
}

func (notConfiguredClient) CancelSubscription(context.Context, string) error {
	return ErrNotConfigured
}

func (notConfiguredClient) ConstructWebhookEvent([]byte, string, string) (stripe.Event, error) {
	return stripe.Event{}, ErrNotConfigured
}

func (notConfiguredClient) CreateExpressAccount(context.Context, string, string) (string, error) {
	return "", ErrNotConfigured
}

func (notConfiguredClient) CreateOnboardingLink(context.Context, string, string, string) (string, error) {
	return "", ErrNotConfigured
}

func (notConfiguredClient) AccountChargesEnabled(context.Context, string) (bool, error) {
	return false, ErrNotConfigured
}

func (notConfiguredClient) ChargeBuyer(context.Context, string, string, int64, string) (string, error) {
	return "", ErrNotConfigured
}

func (notConfiguredClient) ReleaseToSeller(context.Context, string, int64, string, string) (string, error) {
	return "", ErrNotConfigured
}

func strPtr(s string) *string { return &s }
