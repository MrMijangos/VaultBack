package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"vault-payment/src/features/ads/application"
)

type RegisterClickController struct {
	useCase *application.RegisterClickUseCase
}

func NewRegisterClickController(useCase *application.RegisterClickUseCase) *RegisterClickController {
	return &RegisterClickController{useCase: useCase}
}

func (ctrl *RegisterClickController) Handle(c *gin.Context) {
	if err := ctrl.useCase.Execute(c.Request.Context(), c.Param("id")); err != nil {
		c.JSON(statusForError(err), gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
