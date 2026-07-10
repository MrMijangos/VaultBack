package controllers

import (
	"errors"
	"net/http"

	"github.com/google/uuid"

	"vault/src/core/httpresponse"
	"vault/src/features/users/application"
	"vault/src/features/users/domain/repositories"
)

const maxUploadSize = 5 << 20

type UploadUserImageController struct {
	useCase *application.UploadUserImageUseCase
}

func NewUploadUserImageController(useCase *application.UploadUserImageUseCase) *UploadUserImageController {
	return &UploadUserImageController{useCase: useCase}
}

func (c *UploadUserImageController) Handle(w http.ResponseWriter, r *http.Request) {
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

	updated, err := c.useCase.Execute(r.Context(), id, file)
	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			httpresponse.WriteError(w, http.StatusNotFound, err.Error())
			return
		}
		httpresponse.WriteError(w, http.StatusBadGateway, err.Error())
		return
	}

	httpresponse.WriteJSON(w, http.StatusOK, updated)
}
