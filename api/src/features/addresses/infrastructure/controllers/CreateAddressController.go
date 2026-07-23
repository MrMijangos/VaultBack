package controllers

import (
	"encoding/json"
	"net/http"

	"vault/src/core/httpresponse"
	"vault/src/core/security"
	"vault/src/features/addresses/application"
	"vault/src/features/addresses/domain/dto/request"
)

type CreateAddressController struct {
	useCase *application.CreateAddressUseCase
}

func NewCreateAddressController(useCase *application.CreateAddressUseCase) *CreateAddressController {
	return &CreateAddressController{useCase: useCase}
}

func (c *CreateAddressController) Handle(w http.ResponseWriter, r *http.Request) {
	claims, ok := security.ClaimsFromContext(r.Context())
	if !ok {
		httpresponse.WriteError(w, http.StatusUnauthorized, "no autenticado")
		return
	}

	var req request.CreateAddressRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, "cuerpo de la peticion invalido")
		return
	}

	created, err := c.useCase.Execute(r.Context(), claims.UserID, req)
	if err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	httpresponse.WriteJSON(w, http.StatusCreated, created)
}
