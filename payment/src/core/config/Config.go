package config

type Config struct {
	AppPort             string
	JWTSecret           string
	CORSOrigin          string
	RabbitMQURL         string
	StripeSecretKey     string
	StripeWebhookSecret string
	StripePriceBasico   string
	StripePricePro      string
	StripePricePremium  string
	DBHost              string
	DBPort              string
	DBUser              string
	DBPassword          string
	DBName              string
	DBSSL               string
}

// StripeConfigured indica si hay claves reales de Stripe cargadas. Mientras
// no las haya (cuenta de Stripe todavía no creada), el servicio arranca
// igual -- /subscriptions/plans y /ads/* funcionan normal, pero crear o
// cancelar una suscripción responde con un error claro en vez de fallar en
// el arranque.
func (c *Config) StripeConfigured() bool {
	return c.StripeSecretKey != ""
}
