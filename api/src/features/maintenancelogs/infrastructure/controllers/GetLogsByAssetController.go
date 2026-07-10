package controllers

import (
	"net/http"

	"vault/src/core/httpresponse"
	"vault/src/features/maintenancelogs/application"
)

type GetLogsByAssetController struct {
	useCase *application.GetLogsByAssetUseCase
}

func NewGetLogsByAssetController(useCase *application.GetLogsByAssetUseCase) *GetLogsByAssetController {
	return &GetLogsByAssetController{useCase: useCase}
}

func (c *GetLogsByAssetController) Handle(w http.ResponseWriter, r *http.Request) {
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
