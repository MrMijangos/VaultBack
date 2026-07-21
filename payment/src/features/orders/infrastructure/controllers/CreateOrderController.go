package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"vault-payment/src/core/security"
	"vault-payment/src/features/orders/application"
	"vault-payment/src/features/orders/domain/dto/request"
)

type CreateOrderController struct {
	useCase *application.CreateOrderUseCase
}

func NewCreateOrderController(useCase *application.CreateOrderUseCase) *CreateOrderController {
	return &CreateOrderController{useCase: useCase}
}

func (ctrl *CreateOrderController) Handle(c *gin.Context) {
	claims, ok := security.ClaimsFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no autenticado"})
		return
	}

	var req request.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cuerpo de la petición inválido"})
		return
	}

	order, err := ctrl.useCase.Execute(c.Request.Context(), claims.UserID, req)
	if err != nil {
		c.JSON(statusForError(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, order)
}
