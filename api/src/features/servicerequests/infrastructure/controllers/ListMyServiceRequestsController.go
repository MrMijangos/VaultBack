package controllers

import (
	"net/http"

	"vault/src/core/httpresponse"
	"vault/src/core/security"
	"vault/src/features/servicerequests/application"
)

type ListMyServiceRequestsController struct {
	useCase *application.ListMyServiceRequestsUseCase
}

func NewListMyServiceRequestsController(useCase *application.ListMyServiceRequestsUseCase) *ListMyServiceRequestsController {
	return &ListMyServiceRequestsController{useCase: useCase}
}

func (c *ListMyServiceRequestsController) Handle(w http.ResponseWriter, r *http.Request) {
	claims, ok := security.ClaimsFromContext(r.Context())
	if !ok {
		httpresponse.WriteError(w, http.StatusUnauthorized, "no autenticado")
		return
	}

	list, err := c.useCase.Execute(r.Context(), claims.UserID)
	if err != nil {
		httpresponse.WriteError(w, statusForError(err), err.Error())
		return
	}

	httpresponse.WriteJSON(w, http.StatusOK, list)
}
