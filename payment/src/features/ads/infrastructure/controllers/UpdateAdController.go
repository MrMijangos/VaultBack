package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"vault-payment/src/core/security"
	"vault-payment/src/features/ads/application"
	"vault-payment/src/features/ads/domain/dto/request"
)

type UpdateAdController struct {
	useCase *application.UpdateAdUseCase
}

func NewUpdateAdController(useCase *application.UpdateAdUseCase) *UpdateAdController {
	return &UpdateAdController{useCase: useCase}
}

func (ctrl *UpdateAdController) Handle(c *gin.Context) {
	claims, ok := security.ClaimsFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no autenticado"})
		return
	}

	var req request.UpdateAdRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cuerpo de la petición inválido"})
		return
	}

	ad, err := ctrl.useCase.Execute(c.Request.Context(), claims.UserID, c.Param("id"), req)
	if err != nil {
		c.JSON(statusForError(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ad)
}
