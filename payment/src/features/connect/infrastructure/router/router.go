package router

import (
	"github.com/gin-gonic/gin"

	"vault-payment/src/core/security"
	"vault-payment/src/features/connect/infrastructure/controllers"
)

func RegisterRoutes(
	rg *gin.RouterGroup,
	jwtSecret string,
	createOnboardingLink *controllers.CreateOnboardingLinkController,
	getStatus *controllers.GetAccountStatusController,
) {
	connect := rg.Group("/connect")
	connect.Use(security.RequireAuth(jwtSecret))

	connect.POST("/onboarding", createOnboardingLink.Handle)
	connect.GET("/status", getStatus.Handle)
}
