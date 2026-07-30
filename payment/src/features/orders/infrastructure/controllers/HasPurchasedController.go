package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"vault-payment/src/features/orders/application"
)

type HasPurchasedController struct {
	useCase *application.HasPurchasedUseCase
}

func NewHasPurchasedController(useCase *application.HasPurchasedUseCase) *HasPurchasedController {
	return &HasPurchasedController{useCase: useCase}
}

// Handle atiende GET /api/v1/orders/has-purchased?buyer_id=&seller_id= --
// ruta pública (sin JWT) para que reviews/ en api/ pueda consultarla
// servidor a servidor, igual que api/ expone GET /assets/{id} para que este
// mismo servicio consulte el precio real de un activo.
func (ctrl *HasPurchasedController) Handle(c *gin.Context) {
	buyerID := c.Query("buyer_id")
	sellerID := c.Query("seller_id")

	if _, err := uuid.Parse(buyerID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "buyer_id invalido"})
		return
	}
	if _, err := uuid.Parse(sellerID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "seller_id invalido"})
		return
	}

	purchased, err := ctrl.useCase.Execute(c.Request.Context(), buyerID, sellerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"purchased": purchased})
}
