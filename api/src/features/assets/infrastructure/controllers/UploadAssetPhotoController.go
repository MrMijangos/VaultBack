package controllers

import (
	"errors"
	"net/http"

	"github.com/google/uuid"

	"vault/src/core/httpresponse"
	"vault/src/core/security"
	"vault/src/features/assets/application"
	"vault/src/features/assets/domain/repositories"
)

const maxUploadSize = 5 << 20

type UploadAssetPhotoController struct {
	useCase *application.UploadAssetPhotoUseCase
}

func NewUploadAssetPhotoController(useCase *application.UploadAssetPhotoUseCase) *UploadAssetPhotoController {
	return &UploadAssetPhotoController{useCase: useCase}
}

func (c *UploadAssetPhotoController) Handle(w http.ResponseWriter, r *http.Request) {
	claims, ok := security.ClaimsFromContext(r.Context())
	if !ok {
		httpresponse.WriteError(w, http.StatusUnauthorized, "no autenticado")
		return
	}

	id := r.PathValue("id")
	if _, err := uuid.Parse(id); err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, "id invalido")
		return
	}

	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, "la imagen supera el tamaño maximo permitido (5MB)")
		return
	}

	file, _, err := r.FormFile("image")
	if err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, "falta el archivo 'image' en el formulario")
		return
	}
	defer file.Close()

	updated, err := c.useCase.Execute(r.Context(), id, claims.UserID, file)
	if err != nil {
		if errors.Is(err, repositories.ErrAssetNotFound) {
			httpresponse.WriteError(w, http.StatusNotFound, err.Error())
			return
		}
		httpresponse.WriteError(w, http.StatusBadGateway, err.Error())
		return
	}

	httpresponse.WriteJSON(w, http.StatusOK, updated)
}
