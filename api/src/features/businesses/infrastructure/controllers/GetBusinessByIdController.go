package controllers

import (
	"errors"
	"net/http"

	"github.com/google/uuid"

	"vault/src/core/httpresponse"
	"vault/src/features/businesses/application"
	"vault/src/features/businesses/domain/repositories"
)

type GetBusinessByIdController struct {
	useCase *application.GetBusinessByIdUseCase
}

func NewGetBusinessByIdController(useCase *application.GetBusinessByIdUseCase) *GetBusinessByIdController {
	return &GetBusinessByIdController{useCase: useCase}
}

func (c *GetBusinessByIdController) Handle(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if _, err := uuid.Parse(id); err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, "id invalido")
		return
	}

	b, err := c.useCase.Execute(r.Context(), id)
	if err != nil {
		if errors.Is(err, repositories.ErrBusinessNotFound) {
			httpresponse.WriteError(w, http.StatusNotFound, err.Error())
			return
		}
		httpresponse.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	httpresponse.WriteJSON(w, http.StatusOK, b)
}
