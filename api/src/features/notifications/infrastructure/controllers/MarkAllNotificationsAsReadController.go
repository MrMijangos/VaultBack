package controllers

import (
	"net/http"

	"vault/src/core/httpresponse"
	"vault/src/core/security"
	"vault/src/features/notifications/application"
)

type MarkAllNotificationsAsReadController struct {
	useCase *application.MarkAllNotificationsAsReadUseCase
}

func NewMarkAllNotificationsAsReadController(useCase *application.MarkAllNotificationsAsReadUseCase) *MarkAllNotificationsAsReadController {
	return &MarkAllNotificationsAsReadController{useCase: useCase}
}

func (c *MarkAllNotificationsAsReadController) Handle(w http.ResponseWriter, r *http.Request) {
	claims, ok := security.ClaimsFromContext(r.Context())
	if !ok {
		httpresponse.WriteError(w, http.StatusUnauthorized, "no autenticado")
		return
	}

	if err := c.useCase.Execute(r.Context(), claims.UserID); err != nil {
		httpresponse.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
