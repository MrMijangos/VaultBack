package controllers

import (
	"errors"
	"net/http"

	"github.com/google/uuid"

	"vault/src/core/httpresponse"
	"vault/src/features/users/application"
	"vault/src/features/users/domain/repositories"
)

type DeleteUserController struct {
	useCase *application.DeleteUserUseCase
}

func NewDeleteUserController(useCase *application.DeleteUserUseCase) *DeleteUserController {
	return &DeleteUserController{useCase: useCase}
}

func (c *DeleteUserController) Handle(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if _, err := uuid.Parse(id); err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, "id invalido")
		return
	}

	if err := c.useCase.Execute(r.Context(), id); err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			httpresponse.WriteError(w, http.StatusNotFound, err.Error())
			return
		}
		httpresponse.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
