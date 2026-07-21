package router

import (
	"github.com/gin-gonic/gin"

	"vault-payment/src/core/security"
	"vault-payment/src/features/orders/infrastructure/controllers"
)

func RegisterRoutes(
	rg *gin.RouterGroup,
	jwtSecret string,
	create *controllers.CreateOrderController,
	confirm *controllers.ConfirmOrderController,
	get *controllers.GetOrderController,
) {
	orders := rg.Group("/orders")
	orders.Use(security.RequireAuth(jwtSecret))

	orders.POST("", create.Handle)
	orders.GET("/:id", get.Handle)
	orders.POST("/:id/confirm", confirm.Handle)
}
