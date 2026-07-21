package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"vault-payment/src/core/security"
	"vault-payment/src/features/subscriptions/application"
)

type CancelSubscriptionController struct {
	useCase *application.CancelSubscriptionUseCase
}

func NewCancelSubscriptionController(useCase *application.CancelSubscriptionUseCase) *CancelSubscriptionController {
	return &CancelSubscriptionController{useCase: useCase}
}

func (ctrl *CancelSubscriptionController) Handle(c *gin.Context) {
	claims, ok := security.ClaimsFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no autenticado"})
		return
	}

	if err := ctrl.useCase.Execute(c.Request.Context(), claims.UserID); err != nil {
		c.JSON(statusForError(err), gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
