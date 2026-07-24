package controllers

import (
	"net/http"

	"vault/src/core/httpresponse"
	"vault/src/core/security"
	"vault/src/features/assets/application"
)

type GetMyAssetsController struct {
	useCase *application.GetMyAssetsUseCase
}

func NewGetMyAssetsController(useCase *application.GetMyAssetsUseCase) *GetMyAssetsController {
	return &GetMyAssetsController{useCase: useCase}
}

func (c *GetMyAssetsController) Handle(w http.ResponseWriter, r *http.Request) {
	claims, ok := security.ClaimsFromContext(r.Context())
	if !ok {
		httpresponse.WriteError(w, http.StatusUnauthorized, "no autenticado")
		return
	}

	assets, err := c.useCase.Execute(r.Context(), claims.UserID)
	if err != nil {
		httpresponse.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httpresponse.WriteJSON(w, http.StatusOK, assets)
}
