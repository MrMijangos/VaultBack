package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"vault-payment/src/features/ads/application"
)

type ListActiveAdsController struct {
	useCase *application.ListActiveAdsUseCase
}

func NewListActiveAdsController(useCase *application.ListActiveAdsUseCase) *ListActiveAdsController {
	return &ListActiveAdsController{useCase: useCase}
}

func (ctrl *ListActiveAdsController) Handle(c *gin.Context) {
	section := c.Query("section")
	if section == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "el parámetro section es obligatorio"})
		return
	}

	ads, err := ctrl.useCase.Execute(c.Request.Context(), section)
	if err != nil {
		c.JSON(statusForError(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ads": ads})
}
