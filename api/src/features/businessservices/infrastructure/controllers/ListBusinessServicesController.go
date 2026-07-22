package controllers

import (
	"net/http"

	"github.com/google/uuid"

	"vault/src/core/httpresponse"
	"vault/src/features/businessservices/application"
)

type ListBusinessServicesController struct {
	useCase *application.ListBusinessServicesUseCase
}

func NewListBusinessServicesController(useCase *application.ListBusinessServicesUseCase) *ListBusinessServicesController {
	return &ListBusinessServicesController{useCase: useCase}
}

func (c *ListBusinessServicesController) Handle(w http.ResponseWriter, r *http.Request) {
	businessID := r.PathValue("id")
	if _, err := uuid.Parse(businessID); err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, "id de negocio invalido")
		return
	}

	list, err := c.useCase.Execute(r.Context(), businessID)
	if err != nil {
		httpresponse.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	httpresponse.WriteJSON(w, http.StatusOK, list)
}
