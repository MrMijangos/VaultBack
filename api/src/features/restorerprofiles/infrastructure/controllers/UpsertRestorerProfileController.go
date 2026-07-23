package controllers

import (
	"encoding/json"
	"net/http"

	"vault/src/core/httpresponse"
	"vault/src/core/security"
	"vault/src/features/restorerprofiles/application"
	"vault/src/features/restorerprofiles/domain/dto/request"
)

type UpsertRestorerProfileController struct {
	useCase *application.UpsertRestorerProfileUseCase
}

func NewUpsertRestorerProfileController(useCase *application.UpsertRestorerProfileUseCase) *UpsertRestorerProfileController {
	return &UpsertRestorerProfileController{useCase: useCase}
}

func (c *UpsertRestorerProfileController) Handle(w http.ResponseWriter, r *http.Request) {
	claims, ok := security.ClaimsFromContext(r.Context())
	if !ok {
		httpresponse.WriteError(w, http.StatusUnauthorized, "no autenticado")
		return
	}

	userID := r.PathValue("userId")
	if claims.UserID != userID {
		httpresponse.WriteError(w, http.StatusForbidden, "no puedes editar el perfil de otro usuario")
		return
	}

	var req request.UpsertRestorerProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, "cuerpo de la peticion invalido")
		return
	}

	updated, err := c.useCase.Execute(r.Context(), userID, req)
	if err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	httpresponse.WriteJSON(w, http.StatusOK, updated)
}
