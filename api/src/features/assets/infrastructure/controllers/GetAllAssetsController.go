package controllers

import (
	"net/http"

	"vault/src/core/httpresponse"
	"vault/src/features/assets/application"
)

type GetAllAssetsController struct {
	useCase *application.GetAllAssetsUseCase
}

func NewGetAllAssetsController(useCase *application.GetAllAssetsUseCase) *GetAllAssetsController {
	return &GetAllAssetsController{useCase: useCase}
}

func (c *GetAllAssetsController) Handle(w http.ResponseWriter, r *http.Request) {
	assets, err := c.useCase.Execute(r.Context())
	if err != nil {
		httpresponse.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httpresponse.WriteJSON(w, http.StatusOK, assets)
}
