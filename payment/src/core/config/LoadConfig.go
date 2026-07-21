package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		AppPort:             os.Getenv("APP_PORT"),
		JWTSecret:           os.Getenv("JWT_SECRET"),
		CORSOrigin:          os.Getenv("CORS_ORIGIN"),
		RabbitMQURL:         os.Getenv("RABBITMQ_URL"),
		StripeSecretKey:     os.Getenv("STRIPE_SECRET_KEY"),
		StripeWebhookSecret: os.Getenv("STRIPE_WEBHOOK_SECRET"),
		StripePriceBasico:   os.Getenv("STRIPE_PRICE_BASICO"),
		StripePricePro:      os.Getenv("STRIPE_PRICE_PRO"),
		StripePricePremium:  os.Getenv("STRIPE_PRICE_PREMIUM"),
	}

	if cfg.AppPort == "" {
		cfg.AppPort = os.Getenv("PORT")
	}
	if cfg.AppPort == "" {
		cfg.AppPort = "8005"
	}
	if cfg.CORSOrigin == "" {
		cfg.CORSOrigin = "*"
	}
	if cfg.RabbitMQURL == "" {
		cfg.RabbitMQURL = "amqp://guest:guest@localhost:5672/"
	}

	// JWT_SECRET es la única variable realmente obligatoria: sin ella no se
	// puede validar ningún token, y debe ser idéntica a la de api/ (mismos
	// tokens, emitidos por api, validados aquí). Las de Stripe se validan
	// en caliente (ver Config.StripeConfigured) porque todavía no existe la
	// cuenta de Stripe -- el servicio debe poder arrancar sin ellas.
	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("falta la variable de entorno obligatoria: JWT_SECRET")
	}

	return cfg, nil
}
