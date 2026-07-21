package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"vault-payment/src/core/security"
	"vault-payment/src/features/orders/application"
)

type GetOrderController struct {
	useCase *application.GetOrderUseCase
}

func NewGetOrderController(useCase *application.GetOrderUseCase) *GetOrderController {
	return &GetOrderController{useCase: useCase}
}

func (ctrl *GetOrderController) Handle(c *gin.Context) {
	claims, ok := security.ClaimsFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no autenticado"})
		return
	}

	order, err := ctrl.useCase.Execute(c.Request.Context(), claims.UserID, c.Param("id"))
	if err != nil {
		c.JSON(statusForError(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, order)
}
