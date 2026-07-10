package controllers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"

	"vault/src/core/httpresponse"
	"vault/src/features/users/application"
	"vault/src/features/users/domain/dto/request"
	"vault/src/features/users/domain/repositories"
)

type UpdateUserController struct {
	useCase *application.UpdateUserUseCase
}

func NewUpdateUserController(useCase *application.UpdateUserUseCase) *UpdateUserController {
	return &UpdateUserController{useCase: useCase}
}

func (c *UpdateUserController) Handle(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if _, err := uuid.Parse(id); err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, "id invalido")
		return
	}

	var req request.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, "cuerpo de la peticion invalido")
		return
	}

	updated, err := c.useCase.Execute(r.Context(), id, req)
	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			httpresponse.WriteError(w, http.StatusNotFound, err.Error())
			return
		}
		httpresponse.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	httpresponse.WriteJSON(w, http.StatusOK, updated)
}
