package controllers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"

	"vault/src/core/httpresponse"
	"vault/src/core/security"
	"vault/src/features/businesses/application"
	"vault/src/features/businesses/domain/dto/request"
	"vault/src/features/businesses/domain/repositories"
)

type UpdateBusinessController struct {
	useCase *application.UpdateBusinessUseCase
}

func NewUpdateBusinessController(useCase *application.UpdateBusinessUseCase) *UpdateBusinessController {
	return &UpdateBusinessController{useCase: useCase}
}

func (c *UpdateBusinessController) Handle(w http.ResponseWriter, r *http.Request) {
	claims, ok := security.ClaimsFromContext(r.Context())
	if !ok {
		httpresponse.WriteError(w, http.StatusUnauthorized, "no autenticado")
		return
	}

	id := r.PathValue("id")
	if _, err := uuid.Parse(id); err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, "id invalido")
		return
	}

	var req request.UpdateBusinessRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, "cuerpo de la peticion invalido")
		return
	}

	updated, err := c.useCase.Execute(r.Context(), id, claims.UserID, req)
	if err != nil {
		if errors.Is(err, repositories.ErrBusinessNotFound) {
			httpresponse.WriteError(w, http.StatusNotFound, err.Error())
			return
		}
		httpresponse.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	httpresponse.WriteJSON(w, http.StatusOK, updated)
}
