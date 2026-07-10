package router

import (
	"net/http"

	"vault/src/core/security"
	"vault/src/features/blockchaincertificates/infrastructure/controllers"
)

func RegisterRoutes(
	mux *http.ServeMux,
	createCertificate *controllers.CreateBlockchainCertificateController,
	getCertificatesByAsset *controllers.GetCertificatesByAssetController,
	getCertificateById *controllers.GetCertificateByIdController,
	jwtSecret string,
) {
	auth := security.RequireAuth(jwtSecret)

	mux.Handle("POST /api/v1/blockchain-certificates", auth(http.HandlerFunc(createCertificate.Handle)))
	mux.HandleFunc("GET /api/v1/blockchain-certificates", getCertificatesByAsset.Handle)
	mux.HandleFunc("GET /api/v1/blockchain-certificates/{id}", getCertificateById.Handle)
}
