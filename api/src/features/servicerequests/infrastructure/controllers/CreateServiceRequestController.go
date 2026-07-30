package controllers

import (
	"encoding/json"
	"net/http"

	"vault/src/core/httpresponse"
	"vault/src/core/security"
	"vault/src/features/servicerequests/application"
	"vault/src/features/servicerequests/domain/dto/request"
)

type CreateServiceRequestController struct {
	useCase *application.CreateServiceRequestUseCase
}

func NewCreateServiceRequestController(useCase *application.CreateServiceRequestUseCase) *CreateServiceRequestController {
	return &CreateServiceRequestController{useCase: useCase}
}

func (c *CreateServiceRequestController) Handle(w http.ResponseWriter, r *http.Request) {
	claims, ok := security.ClaimsFromContext(r.Context())
	if !ok {
		httpresponse.WriteError(w, http.StatusUnauthorized, "no autenticado")
		return
	}

	var req request.CreateServiceRequestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, "cuerpo de la peticion invalido")
		return
	}

	created, err := c.useCase.Execute(r.Context(), claims.UserID, req)
	if err != nil {
		httpresponse.WriteError(w, statusForError(err), err.Error())
		return
	}

	httpresponse.WriteJSON(w, http.StatusCreated, created)
}
