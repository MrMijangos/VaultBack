package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"vault-payment/src/core/security"
	"vault-payment/src/features/ads/application"
)

type ListMyAdsController struct {
	useCase *application.ListMyAdsUseCase
}

func NewListMyAdsController(useCase *application.ListMyAdsUseCase) *ListMyAdsController {
	return &ListMyAdsController{useCase: useCase}
}

func (ctrl *ListMyAdsController) Handle(c *gin.Context) {
	claims, ok := security.ClaimsFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no autenticado"})
		return
	}

	ads, err := ctrl.useCase.Execute(c.Request.Context(), claims.UserID)
	if err != nil {
		c.JSON(statusForError(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ads": ads})
}
