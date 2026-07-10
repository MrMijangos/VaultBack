package controllers

import (
	"net/http"

	"vault/src/core/httpresponse"
	"vault/src/features/blockchaincertificates/application"
)

type GetCertificatesByAssetController struct {
	useCase *application.GetCertificatesByAssetUseCase
}

func NewGetCertificatesByAssetController(useCase *application.GetCertificatesByAssetUseCase) *GetCertificatesByAssetController {
	return &GetCertificatesByAssetController{useCase: useCase}
}

func (c *GetCertificatesByAssetController) Handle(w http.ResponseWriter, r *http.Request) {
	assetID := r.URL.Query().Get("asset_id")
	if assetID == "" {
		httpresponse.WriteError(w, http.StatusBadRequest, "el parametro asset_id es obligatorio")
		return
	}

	list, err := c.useCase.Execute(r.Context(), assetID)
	if err != nil {
		httpresponse.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httpresponse.WriteJSON(w, http.StatusOK, list)
}
