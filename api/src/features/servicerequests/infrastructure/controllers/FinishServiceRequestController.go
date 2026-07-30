package controllers

import (
	"net/http"

	"vault/src/core/httpresponse"
	"vault/src/core/security"
	"vault/src/features/servicerequests/application"
)

type FinishServiceRequestController struct {
	useCase *application.FinishServiceRequestUseCase
}

func NewFinishServiceRequestController(useCase *application.FinishServiceRequestUseCase) *FinishServiceRequestController {
	return &FinishServiceRequestController{useCase: useCase}
}

func (c *FinishServiceRequestController) Handle(w http.ResponseWriter, r *http.Request) {
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
