package controllers

import (
	"errors"
	"net/http"

	"github.com/google/uuid"

	"vault/src/core/httpresponse"
	"vault/src/core/security"
	"vault/src/features/businesses/application"
	"vault/src/features/businesses/domain/repositories"
)

type DeleteBusinessController struct {
	useCase *application.DeleteBusinessUseCase
}

func NewDeleteBusinessController(useCase *application.DeleteBusinessUseCase) *DeleteBusinessController {
	return &DeleteBusinessController{useCase: useCase}
}

func (c *DeleteBusinessController) Handle(w http.ResponseWriter, r *http.Request) {
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

	if err := c.useCase.Execute(r.Context(), id, claims.UserID); err != nil {
		if errors.Is(err, repositories.ErrBusinessNotFound) {
			httpresponse.WriteError(w, http.StatusNotFound, err.Error())
			return
		}
		httpresponse.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
