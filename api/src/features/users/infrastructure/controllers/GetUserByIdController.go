package controllers

import (
	"errors"
	"net/http"

	"github.com/google/uuid"

	"vault/src/core/httpresponse"
	"vault/src/features/users/application"
	"vault/src/features/users/domain/repositories"
)

type GetUserByIdController struct {
	useCase *application.GetUserByIdUseCase
}

func NewGetUserByIdController(useCase *application.GetUserByIdUseCase) *GetUserByIdController {
	return &GetUserByIdController{useCase: useCase}
}

func (c *GetUserByIdController) Handle(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if _, err := uuid.Parse(id); err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, "id invalido")
		return
	}

	user, err := c.useCase.Execute(r.Context(), id)
	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			httpresponse.WriteError(w, http.StatusNotFound, err.Error())
			return
		}
		httpresponse.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	httpresponse.WriteJSON(w, http.StatusOK, user)
}
