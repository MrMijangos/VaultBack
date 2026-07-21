package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"vault-payment/src/core/security"
	"vault-payment/src/features/orders/application"
)

type ConfirmOrderController struct {
	useCase *application.ConfirmOrderUseCase
}

func NewConfirmOrderController(useCase *application.ConfirmOrderUseCase) *ConfirmOrderController {
	return &ConfirmOrderController{useCase: useCase}
}

func (ctrl *ConfirmOrderController) Handle(c *gin.Context) {
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
