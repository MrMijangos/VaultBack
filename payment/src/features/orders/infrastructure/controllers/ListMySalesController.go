package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"vault-payment/src/core/security"
	"vault-payment/src/features/orders/application"
)

type ListMySalesController struct {
	useCase *application.ListMySalesUseCase
}

func NewListMySalesController(useCase *application.ListMySalesUseCase) *ListMySalesController {
	return &ListMySalesController{useCase: useCase}
}

func (ctrl *ListMySalesController) Handle(c *gin.Context) {
	claims, ok := security.ClaimsFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no autenticado"})
		return
	}

	orders, err := ctrl.useCase.Execute(c.Request.Context(), claims.UserID)
	if err != nil {
		c.JSON(statusForError(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"orders": orders})
}
