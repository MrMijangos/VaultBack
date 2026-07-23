package controllers

import (
	"encoding/json"
	"net/http"

	"vault/src/core/httpresponse"
	"vault/src/core/security"
	"vault/src/features/users/application"
	"vault/src/features/users/domain/dto/request"
)

type SetPublicKeyController struct {
	useCase *application.SetPublicKeyUseCase
}

func NewSetPublicKeyController(useCase *application.SetPublicKeyUseCase) *SetPublicKeyController {
	return &SetPublicKeyController{useCase: useCase}
}

func (c *SetPublicKeyController) Handle(w http.ResponseWriter, r *http.Request) {
	claims, ok := security.ClaimsFromContext(r.Context())
	if !ok {
		httpresponse.WriteError(w, http.StatusUnauthorized, "no autenticado")
		return
	}

	id := r.PathValue("id")
	if claims.UserID != id {
		httpresponse.WriteError(w, http.StatusForbidden, "no puedes registrar la llave publica de otro usuario")
		return
	}

	var req request.SetPublicKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, "cuerpo de la peticion invalido")
		return
	}

	if err := c.useCase.Execute(r.Context(), id, req); err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
