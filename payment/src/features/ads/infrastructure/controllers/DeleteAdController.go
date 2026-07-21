package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"vault-payment/src/core/security"
	"vault-payment/src/features/ads/application"
)

type DeleteAdController struct {
	useCase *application.DeleteAdUseCase
}

func NewDeleteAdController(useCase *application.DeleteAdUseCase) *DeleteAdController {
	return &DeleteAdController{useCase: useCase}
}

func (ctrl *DeleteAdController) Handle(c *gin.Context) {
	claims, ok := security.ClaimsFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no autenticado"})
		return
	}

	if err := ctrl.useCase.Execute(c.Request.Context(), claims.UserID, c.Param("id")); err != nil {
		c.JSON(statusForError(err), gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
