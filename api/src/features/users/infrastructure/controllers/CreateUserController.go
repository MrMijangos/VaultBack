package controllers

import (
	"encoding/json"
	"errors"
	"net/http"

	"vault/src/core/httpresponse"
	"vault/src/features/users/application"
	"vault/src/features/users/domain/dto/request"
)

type CreateUserController struct {
	useCase *application.CreateUserUseCase
}

func NewCreateUserController(useCase *application.CreateUserUseCase) *CreateUserController {
	return &CreateUserController{useCase: useCase}
}

func (c *CreateUserController) Handle(w http.ResponseWriter, r *http.Request) {
	var req request.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, "cuerpo de la peticion invalido")
		return
	}

	created, err := c.useCase.Execute(r.Context(), req)
	if err != nil {
		if errors.Is(err, application.ErrEmailTaken) {
			httpresponse.WriteError(w, http.StatusConflict, err.Error())
			return
		}
		httpresponse.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	httpresponse.WriteJSON(w, http.StatusCreated, created)
}
