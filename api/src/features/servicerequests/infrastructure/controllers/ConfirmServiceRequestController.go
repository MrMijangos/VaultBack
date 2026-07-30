package controllers

import (
	"net/http"

	"vault/src/core/httpresponse"
	"vault/src/core/security"
	"vault/src/features/servicerequests/application"
)

type ConfirmServiceRequestController struct {
	useCase *application.ConfirmServiceRequestUseCase
}

func NewConfirmServiceRequestController(useCase *application.ConfirmServiceRequestUseCase) *ConfirmServiceRequestController {
	return &ConfirmServiceRequestController{useCase: useCase}
}

func (c *ConfirmServiceRequestController) Handle(w http.ResponseWriter, r *http.Request) {
	claims, ok := security.ClaimsFromContext(r.Context())
	if !ok {
		httpresponse.WriteError(w, http.StatusUnauthorized, "no autenticado")
		return
	}

	updated, err := c.useCase.Execute(r.Context(), claims.UserID, r.PathValue("id"))
	if err != nil {
		httpresponse.WriteError(w, statusForError(err), err.Error())
		return
	}

	httpresponse.WriteJSON(w, http.StatusOK, updated)
}
