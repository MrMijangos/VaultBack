package controllers

import (
	"errors"
	"net/http"

	"github.com/google/uuid"

	"vault/src/core/httpresponse"
	"vault/src/features/assets/application"
	"vault/src/features/assets/domain/repositories"
)

type GetAssetByIdController struct {
	useCase *application.GetAssetByIdUseCase
}

func NewGetAssetByIdController(useCase *application.GetAssetByIdUseCase) *GetAssetByIdController {
	return &GetAssetByIdController{useCase: useCase}
}

func (c *GetAssetByIdController) Handle(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if _, err := uuid.Parse(id); err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, "id invalido")
		return
	}

	asset, err := c.useCase.Execute(r.Context(), id)
	if err != nil {
		if errors.Is(err, repositories.ErrAssetNotFound) {
			httpresponse.WriteError(w, http.StatusNotFound, err.Error())
			return
		}
		httpresponse.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	httpresponse.WriteJSON(w, http.StatusOK, asset)
}
