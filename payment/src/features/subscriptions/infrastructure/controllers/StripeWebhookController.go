package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"vault-payment/src/features/subscriptions/application"
)

type StripeWebhookController struct {
	useCase *application.HandleStripeWebhookUseCase
}

func NewStripeWebhookController(useCase *application.HandleStripeWebhookUseCase) *StripeWebhookController {
	return &StripeWebhookController{useCase: useCase}
}

func (ctrl *StripeWebhookController) Handle(c *gin.Context) {
	payload, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no se pudo leer el cuerpo de la petición"})
		return
	}

	sigHeader := c.GetHeader("Stripe-Signature")

	if err := ctrl.useCase.Execute(c.Request.Context(), payload, sigHeader); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
