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
	ship *controllers.ShipOrderController,
	listMine *controllers.ListMyOrdersController,
	listSelling *controllers.ListMySalesController,
	hasPurchased *controllers.HasPurchasedController,
) {
	orders := rg.Group("/orders")

	// Pública (sin RequireAuth): la consulta reviews/ en api/ servidor a
	// servidor para exigir una compra completada antes de dejar reseñar a
	// un vendedor -- no lleva el JWT del comprador, solo los IDs.
	orders.GET("/has-purchased", hasPurchased.Handle)

	orders.Use(security.RequireAuth(jwtSecret))

	orders.POST("", create.Handle)
	orders.GET("/mine", listMine.Handle)
	orders.GET("/selling", listSelling.Handle)
	orders.GET("/:id", get.Handle)
	orders.POST("/:id/confirm", confirm.Handle)
	orders.POST("/:id/ship", ship.Handle)
}
