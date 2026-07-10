package controllers

import (
	"net/http"

	"vault/src/core/httpresponse"
	"vault/src/features/users/application"
)

type GetAllUsersController struct {
	useCase *application.GetAllUsersUseCase
}

func NewGetAllUsersController(useCase *application.GetAllUsersUseCase) *GetAllUsersController {
	return &GetAllUsersController{useCase: useCase}
}

func (c *GetAllUsersController) Handle(w http.ResponseWriter, r *http.Request) {
	users, err := c.useCase.Execute(r.Context())
	if err != nil {
		httpresponse.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httpresponse.WriteJSON(w, http.StatusOK, users)
}
