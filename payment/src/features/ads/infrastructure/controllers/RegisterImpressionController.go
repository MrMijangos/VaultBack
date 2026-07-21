package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"vault-payment/src/features/ads/application"
)

type RegisterImpressionController struct {
	useCase *application.RegisterImpressionUseCase
}

func NewRegisterImpressionController(useCase *application.RegisterImpressionUseCase) *RegisterImpressionController {
	return &RegisterImpressionController{useCase: useCase}
}

func (ctrl *RegisterImpressionController) Handle(c *gin.Context) {
	if err := ctrl.useCase.Execute(c.Request.Context(), c.Param("id")); err != nil {
		c.JSON(statusForError(err), gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
