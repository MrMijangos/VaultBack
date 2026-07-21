package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"vault-payment/src/core/security"
	"vault-payment/src/core/stripeclient"
	"vault-payment/src/features/connect/application"
)

type GetAccountStatusController struct {
	useCase *application.GetAccountStatusUseCase
}

func NewGetAccountStatusController(useCase *application.GetAccountStatusUseCase) *GetAccountStatusController {
	return &GetAccountStatusController{useCase: useCase}
}

func (ctrl *GetAccountStatusController) Handle(c *gin.Context) {
	claims, ok := security.ClaimsFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no autenticado"})
		return
	}

	status, err := ctrl.useCase.Execute(c.Request.Context(), claims.UserID)
	if err != nil {
		httpStatus := http.StatusInternalServerError
		if stripeclient.IsNotConfigured(err) {
			httpStatus = http.StatusServiceUnavailable
		}
		c.JSON(httpStatus, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"account": status})
}
