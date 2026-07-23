package controllers

import (
	"errors"
	"net/http"

	"vault/src/core/httpresponse"
	"vault/src/features/users/application"
	"vault/src/features/users/domain/repositories"
)

type GetPublicKeyController struct {
	useCase *application.GetPublicKeyUseCase
}

func NewGetPublicKeyController(useCase *application.GetPublicKeyUseCase) *GetPublicKeyController {
	return &GetPublicKeyController{useCase: useCase}
}

func (c *GetPublicKeyController) Handle(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	result, err := c.useCase.Execute(r.Context(), id)
	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			httpresponse.WriteError(w, http.StatusNotFound, err.Error())
			return
		}
		httpresponse.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	httpresponse.WriteJSON(w, http.StatusOK, result)
}
