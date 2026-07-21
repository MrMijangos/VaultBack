package router

import (
	"github.com/gin-gonic/gin"

	"vault-payment/src/core/security"
	"vault-payment/src/features/subscriptions/infrastructure/controllers"
)

// RegisterRoutes registra las rutas de suscripciones. El webhook de Stripe
// no lleva RequireAuth: Stripe no manda un JWT, la petición se autentica
// verificando la firma (Stripe-Signature) dentro del propio use case.
func RegisterRoutes(
	rg *gin.RouterGroup,
	jwtSecret string,
	listPlans *controllers.ListPlansController,
	createSubscription *controllers.CreateSubscriptionController,
	getStatus *controllers.GetSubscriptionStatusController,
	cancelSubscription *controllers.CancelSubscriptionController,
	stripeWebhook *controllers.StripeWebhookController,
) {
	sub := rg.Group("/subscriptions")

	sub.GET("/plans", listPlans.Handle)
	sub.POST("/webhook", stripeWebhook.Handle)

	authed := sub.Group("")
	authed.Use(security.RequireAuth(jwtSecret))
	authed.POST("", createSubscription.Handle)
	authed.GET("/me", getStatus.Handle)
	authed.DELETE("", cancelSubscription.Handle)
}
