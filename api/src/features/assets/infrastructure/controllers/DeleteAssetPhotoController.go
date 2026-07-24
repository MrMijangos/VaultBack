package controllers

import (
	"errors"
	"net/http"

	"vault/src/core/httpresponse"
	"vault/src/core/security"
	"vault/src/features/assets/application"
	"vault/src/features/assets/domain/repositories"
)

type DeleteAssetPhotoController struct {
	useCase *application.DeleteAssetPhotoUseCase
}

func NewDeleteAssetPhotoController(useCase *application.DeleteAssetPhotoUseCase) *DeleteAssetPhotoController {
	return &DeleteAssetPhotoController{useCase: useCase}
}

func (c *DeleteAssetPhotoController) Handle(w http.ResponseWriter, r *http.Request) {
	claims, ok := security.ClaimsFromContext(r.Context())
	if !ok {
		httpresponse.WriteError(w, http.StatusUnauthorized, "no autenticado")
		return
	}

	assetID := r.PathValue("id")
	photoID := r.PathValue("photoId")

	updated, err := c.useCase.Execute(r.Context(), assetID, photoID, claims.UserID)
	if err != nil {
		if errors.Is(err, repositories.ErrAssetNotFound) {
			httpresponse.WriteError(w, http.StatusNotFound, err.Error())
			return
		}
		httpresponse.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	httpresponse.WriteJSON(w, http.StatusOK, updated)
}
