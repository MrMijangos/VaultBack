package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"vault-payment/src/core/security"
	"vault-payment/src/features/subscriptions/application"
	"vault-payment/src/features/subscriptions/domain/dto/request"
)

type CreateSubscriptionController struct {
	useCase *application.CreateSubscriptionUseCase
}

func NewCreateSubscriptionController(useCase *application.CreateSubscriptionUseCase) *CreateSubscriptionController {
	return &CreateSubscriptionController{useCase: useCase}
}

func (ctrl *CreateSubscriptionController) Handle(c *gin.Context) {
	claims, ok := security.ClaimsFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no autenticado"})
		return
	}

	var req request.CreateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cuerpo de la petición inválido"})
		return
	}

	sub, err := ctrl.useCase.Execute(c.Request.Context(), claims.UserID, claims.Role, req)
	if err != nil {
		c.JSON(statusForError(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, sub)
}
