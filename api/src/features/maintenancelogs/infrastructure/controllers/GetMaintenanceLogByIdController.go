package controllers

import (
	"errors"
	"net/http"

	"github.com/google/uuid"

	"vault/src/core/httpresponse"
	"vault/src/features/maintenancelogs/application"
	"vault/src/features/maintenancelogs/domain/repositories"
)

type GetMaintenanceLogByIdController struct {
	useCase *application.GetMaintenanceLogByIdUseCase
}

func NewGetMaintenanceLogByIdController(useCase *application.GetMaintenanceLogByIdUseCase) *GetMaintenanceLogByIdController {
	return &GetMaintenanceLogByIdController{useCase: useCase}
}

func (c *GetMaintenanceLogByIdController) Handle(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if _, err := uuid.Parse(id); err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, "id invalido")
		return
	}

	l, err := c.useCase.Execute(r.Context(), id)
	if err != nil {
		if errors.Is(err, repositories.ErrMaintenanceLogNotFound) {
			httpresponse.WriteError(w, http.StatusNotFound, err.Error())
			return
		}
		httpresponse.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	httpresponse.WriteJSON(w, http.StatusOK, l)
}
