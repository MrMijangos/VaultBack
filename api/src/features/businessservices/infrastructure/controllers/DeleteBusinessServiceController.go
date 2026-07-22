package controllers

import (
	"errors"
	"net/http"

	"github.com/google/uuid"

	"vault/src/core/httpresponse"
	"vault/src/core/security"
	"vault/src/features/businessservices/application"
	"vault/src/features/businessservices/domain/repositories"
)

type DeleteBusinessServiceController struct {
	useCase *application.DeleteBusinessServiceUseCase
}

func NewDeleteBusinessServiceController(useCase *application.DeleteBusinessServiceUseCase) *DeleteBusinessServiceController {
	return &DeleteBusinessServiceController{useCase: useCase}
}

func (c *DeleteBusinessServiceController) Handle(w http.ResponseWriter, r *http.Request) {
	claims, ok := security.ClaimsFromContext(r.Context())
	if !ok {
		httpresponse.WriteError(w, http.StatusUnauthorized, "no autenticado")
		return
	}

	businessID := r.PathValue("id")
	if _, err := uuid.Parse(businessID); err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, "id de negocio invalido")
		return
	}

	serviceID := r.PathValue("serviceId")
	if _, err := uuid.Parse(serviceID); err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, "id de servicio invalido")
		return
	}

	if err := c.useCase.Execute(r.Context(), businessID, serviceID, claims.UserID); err != nil {
		switch {
		case errors.Is(err, repositories.ErrBusinessNotFound), errors.Is(err, repositories.ErrBusinessServiceNotFound):
			httpresponse.WriteError(w, http.StatusNotFound, err.Error())
		case errors.Is(err, repositories.ErrNotOwner):
			httpresponse.WriteError(w, http.StatusForbidden, err.Error())
		default:
			httpresponse.WriteError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
