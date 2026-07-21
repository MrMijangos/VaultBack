package router

import (
	"github.com/gin-gonic/gin"

	"vault-payment/src/core/security"
	"vault-payment/src/features/ads/infrastructure/controllers"
)

// RegisterRoutes registra las rutas de anuncios. Listar anuncios activos y
// registrar impressions/clicks son públicas (las consume el feed/marketplace
// de cualquier usuario, no solo el dueño del anuncio); crear/editar/borrar
// requieren sesión porque solo el dueño puede administrar sus anuncios.
func RegisterRoutes(
	rg *gin.RouterGroup,
	jwtSecret string,
	create *controllers.CreateAdController,
	update *controllers.UpdateAdController,
	deleteAd *controllers.DeleteAdController,
	listActive *controllers.ListActiveAdsController,
	registerImpression *controllers.RegisterImpressionController,
	registerClick *controllers.RegisterClickController,
) {
	ads := rg.Group("/ads")

	ads.GET("", listActive.Handle)
	ads.POST("/:id/impression", registerImpression.Handle)
	ads.POST("/:id/click", registerClick.Handle)

	authed := ads.Group("")
	authed.Use(security.RequireAuth(jwtSecret))
	authed.POST("", create.Handle)
	authed.PUT("/:id", update.Handle)
	authed.DELETE("/:id", deleteAd.Handle)
}
