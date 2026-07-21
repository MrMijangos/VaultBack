package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"vault-payment/src/core/security"
	"vault-payment/src/core/stripeclient"
	"vault-payment/src/features/connect/application"
	"vault-payment/src/features/connect/domain/dto/request"
)

type CreateOnboardingLinkController struct {
	useCase *application.CreateOnboardingLinkUseCase
}

func NewCreateOnboardingLinkController(useCase *application.CreateOnboardingLinkUseCase) *CreateOnboardingLinkController {
	return &CreateOnboardingLinkController{useCase: useCase}
}

func (ctrl *CreateOnboardingLinkController) Handle(c *gin.Context) {
	claims, ok := security.ClaimsFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no autenticado"})
		return
	}

	var req request.CreateOnboardingLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.Email == "" || req.RefreshURL == "" || req.ReturnURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email, refresh_url y return_url son obligatorios"})
		return
	}

	url, err := ctrl.useCase.Execute(c.Request.Context(), claims.UserID, req.Email, req.RefreshURL, req.ReturnURL)
	if err != nil {
		status := http.StatusInternalServerError
		if stripeclient.IsNotConfigured(err) {
			status = http.StatusServiceUnavailable
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": url})
}
