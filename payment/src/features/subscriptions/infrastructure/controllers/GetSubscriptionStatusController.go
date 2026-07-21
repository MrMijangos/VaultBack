package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"vault-payment/src/core/security"
	"vault-payment/src/features/subscriptions/application"
)

type GetSubscriptionStatusController struct {
	useCase *application.GetSubscriptionStatusUseCase
}

func NewGetSubscriptionStatusController(useCase *application.GetSubscriptionStatusUseCase) *GetSubscriptionStatusController {
	return &GetSubscriptionStatusController{useCase: useCase}
}

func (ctrl *GetSubscriptionStatusController) Handle(c *gin.Context) {
	claims, ok := security.ClaimsFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no autenticado"})
		return
	}

	sub, err := ctrl.useCase.Execute(c.Request.Context(), claims.UserID)
	if err != nil {
		c.JSON(statusForError(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"subscription": sub})
}
