package controllers

import (
	"errors"
	"net/http"

	"vault/src/core/httpresponse"
	"vault/src/core/security"
	"vault/src/features/businesses/application"
	"vault/src/features/businesses/domain/repositories"
)

type DeleteBusinessPhotoController struct {
	useCase *application.DeleteBusinessPhotoUseCase
}

func NewDeleteBusinessPhotoController(useCase *application.DeleteBusinessPhotoUseCase) *DeleteBusinessPhotoController {
	return &DeleteBusinessPhotoController{useCase: useCase}
}

func (c *DeleteBusinessPhotoController) Handle(w http.ResponseWriter, r *http.Request) {
	claims, ok := security.ClaimsFromContext(r.Context())
	if !ok {
		httpresponse.WriteError(w, http.StatusUnauthorized, "no autenticado")
		return
	}

	businessID := r.PathValue("id")
	photoID := r.PathValue("photoId")

	updated, err := c.useCase.Execute(r.Context(), businessID, photoID, claims.UserID)
	if err != nil {
		if errors.Is(err, repositories.ErrBusinessNotFound) {
			httpresponse.WriteError(w, http.StatusNotFound, err.Error())
			return
		}
		httpresponse.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	httpresponse.WriteJSON(w, http.StatusOK, updated)
}
