package controllers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"

	"vault/src/core/httpresponse"
	"vault/src/core/security"
	"vault/src/features/businessservices/application"
	"vault/src/features/businessservices/domain/dto/request"
	"vault/src/features/businessservices/domain/repositories"
)

type UpdateBusinessServiceController struct {
	useCase *application.UpdateBusinessServiceUseCase
}

func NewUpdateBusinessServiceController(useCase *application.UpdateBusinessServiceUseCase) *UpdateBusinessServiceController {
	return &UpdateBusinessServiceController{useCase: useCase}
}

func (c *UpdateBusinessServiceController) Handle(w http.ResponseWriter, r *http.Request) {
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

	var req request.BusinessServiceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, "cuerpo de la peticion invalido")
		return
	}

	updated, err := c.useCase.Execute(r.Context(), businessID, serviceID, claims.UserID, req)
	if err != nil {
		switch {
		case errors.Is(err, repositories.ErrBusinessNotFound), errors.Is(err, repositories.ErrBusinessServiceNotFound):
			httpresponse.WriteError(w, http.StatusNotFound, err.Error())
		case errors.Is(err, repositories.ErrNotOwner):
			httpresponse.WriteError(w, http.StatusForbidden, err.Error())
		default:
			httpresponse.WriteError(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	httpresponse.WriteJSON(w, http.StatusOK, updated)
}
