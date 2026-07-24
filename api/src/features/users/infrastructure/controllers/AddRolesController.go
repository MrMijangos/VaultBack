package controllers

import (
	"encoding/json"
	"net/http"

	"vault/src/core/httpresponse"
	"vault/src/core/security"
	"vault/src/features/users/application"
	"vault/src/features/users/domain/dto/request"
)

type AddRolesController struct {
	useCase *application.AddRolesUseCase
}

func NewAddRolesController(useCase *application.AddRolesUseCase) *AddRolesController {
	return &AddRolesController{useCase: useCase}
}

func (c *AddRolesController) Handle(w http.ResponseWriter, r *http.Request) {
	claims, ok := security.ClaimsFromContext(r.Context())
	if !ok {
		httpresponse.WriteError(w, http.StatusUnauthorized, "no autenticado")
		return
	}

	id := r.PathValue("id")
	if claims.UserID != id {
		httpresponse.WriteError(w, http.StatusForbidden, "no puedes modificar los roles de otro usuario")
		return
	}

	var req request.AddRolesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, "cuerpo de la peticion invalido")
		return
	}

	updated, err := c.useCase.Execute(r.Context(), id, req)
	if err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	httpresponse.WriteJSON(w, http.StatusOK, updated)
}
