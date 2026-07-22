package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"vault-payment/src/core/config"
	"vault-payment/src/core/eventbus"
	"vault-payment/src/core/stripeclient"
	adsInfra "vault-payment/src/features/ads/infrastructure"
	adsAdapters "vault-payment/src/features/ads/infrastructure/adapters"
	adsRouter "vault-payment/src/features/ads/infrastructure/router"
	connectInfra "vault-payment/src/features/connect/infrastructure"
	connectAdapters "vault-payment/src/features/connect/infrastructure/adapters"
	connectRouter "vault-payment/src/features/connect/infrastructure/router"
	ordersInfra "vault-payment/src/features/orders/infrastructure"
	ordersAdapters "vault-payment/src/features/orders/infrastructure/adapters"
	ordersRouter "vault-payment/src/features/orders/infrastructure/router"
	subscriptionsInfra "vault-payment/src/features/subscriptions/infrastructure"
	subscriptionsAdapters "vault-payment/src/features/subscriptions/infrastructure/adapters"
	subscriptionsRouter "vault-payment/src/features/subscriptions/infrastructure/router"
)

func main() {
	fmt.Println("¡Servicio de pagos de Vault inicializado correctamente!")

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("configuracion invalida: %v", err)
	}

	pool, err := config.NewPoolConnection(cfg)
	if err != nil {
		log.Fatalf("error de conexion a base de datos: %v", err)
	}
	defer pool.Close()

	if err := config.RunMigrations(pool); err != nil {
		log.Fatalf("error al migrar el esquema: %v", err)
	}

	// planRepo se queda en memoria: son 3 planes fijos hardcodeados, no hace
	// falta tabla. El resto ya persiste en Postgres (mismo que usa api/).
	planRepo := subscriptionsAdapters.NewInMemoryPlanRepository(cfg)
	subscriptionRepo := subscriptionsAdapters.NewPostgreSQLSubscriptionRepository(pool)
	adRepo := adsAdapters.NewPostgreSQLAdRepository(pool)
	connectedAccountRepo := connectAdapters.NewPostgreSQLConnectedAccountRepository(pool)
	orderRepo := ordersAdapters.NewPostgreSQLOrderRepository(pool)

	adDeactivator := adsAdapters.NewAdDeactivator(adRepo)
	subscriptionInfoProvider := subscriptionsAdapters.NewSubscriptionInfoAdapter(subscriptionRepo, planRepo)
	sellerCommissionProvider := subscriptionsAdapters.NewSellerCommissionAdapter(subscriptionRepo, planRepo)
	sellerAccountProvider := connectAdapters.NewSellerAccountAdapter(connectedAccountRepo)

	// Sin cuenta de Stripe todavía: el cliente arranca en modo "no
	// configurado" y las rutas que dependen de Stripe responden con un
	// error claro en vez de fallar en el arranque (ver StripeConfigured).
	stripeClient := stripeclient.New(cfg.StripeSecretKey)

	var publisher eventbus.Publisher
	rabbitPublisher, err := eventbus.NewRabbitMQPublisher(cfg.RabbitMQURL)
	if err != nil {
		log.Printf("advertencia: RabbitMQ no disponible (%v), los eventos de suscripción no se publicarán hasta reiniciar el servicio", err)
		publisher = eventbus.NoopPublisher{}
	} else {
		defer rabbitPublisher.Close()
		publisher = rabbitPublisher
	}

	engine := gin.Default()
	engine.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", cfg.CORSOrigin)
		c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":            "ok",
			"stripe_configured": cfg.StripeConfigured(),
		})
	})

	api := engine.Group("/api/v1")

	subscriptionsRouter.RegisterRoutes(
		api,
		cfg.JWTSecret,
		subscriptionsInfra.BuildListPlansController(planRepo),
		subscriptionsInfra.BuildCreateSubscriptionController(planRepo, subscriptionRepo, stripeClient, publisher),
		subscriptionsInfra.BuildGetSubscriptionStatusController(subscriptionRepo),
		subscriptionsInfra.BuildCancelSubscriptionController(subscriptionRepo, stripeClient, adDeactivator, publisher),
		subscriptionsInfra.BuildStripeWebhookController(subscriptionRepo, stripeClient, adDeactivator, publisher, cfg.StripeWebhookSecret),
	)

	adsRouter.RegisterRoutes(
		api,
		cfg.JWTSecret,
		adsInfra.BuildCreateAdController(adRepo, subscriptionInfoProvider),
		adsInfra.BuildUpdateAdController(adRepo),
		adsInfra.BuildDeleteAdController(adRepo),
		adsInfra.BuildListActiveAdsController(adRepo),
		adsInfra.BuildRegisterImpressionController(adRepo),
		adsInfra.BuildRegisterClickController(adRepo),
	)

	connectRouter.RegisterRoutes(
		api,
		cfg.JWTSecret,
		connectInfra.BuildCreateOnboardingLinkController(connectedAccountRepo, stripeClient),
		connectInfra.BuildGetAccountStatusController(connectedAccountRepo, stripeClient),
	)

	ordersRouter.RegisterRoutes(
		api,
		cfg.JWTSecret,
		ordersInfra.BuildCreateOrderController(orderRepo, sellerCommissionProvider, sellerAccountProvider, stripeClient),
		ordersInfra.BuildConfirmOrderController(orderRepo, sellerAccountProvider, stripeClient, publisher),
		ordersInfra.BuildGetOrderController(orderRepo),
	)

	fmt.Println("Servicio de pagos Vault iniciado correctamente.")
	log.Fatal(engine.Run(":" + cfg.AppPort))
}
