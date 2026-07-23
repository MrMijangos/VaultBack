package controllers

import (
	"net/http"

	"github.com/google/uuid"

	"vault/src/core/httpresponse"
	"vault/src/features/assetcomments/application"
)

type GetAssetCommentsController struct {
	useCase *application.GetAssetCommentsUseCase
}

func NewGetAssetCommentsController(useCase *application.GetAssetCommentsUseCase) *GetAssetCommentsController {
	return &GetAssetCommentsController{useCase: useCase}
}

func (c *GetAssetCommentsController) Handle(w http.ResponseWriter, r *http.Request) {
	assetID := r.PathValue("id")
	if _, err := uuid.Parse(assetID); err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, "id invalido")
		return
	}

	list, err := c.useCase.Execute(r.Context(), assetID)
	if err != nil {
		httpresponse.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httpresponse.WriteJSON(w, http.StatusOK, list)
}
