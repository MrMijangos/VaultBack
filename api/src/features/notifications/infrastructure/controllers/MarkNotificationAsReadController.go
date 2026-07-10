package controllers

import (
	"errors"
	"net/http"

	"github.com/google/uuid"

	"vault/src/core/httpresponse"
	"vault/src/core/security"
	"vault/src/features/notifications/application"
	"vault/src/features/notifications/domain/repositories"
)

type MarkNotificationAsReadController struct {
	useCase *application.MarkNotificationAsReadUseCase
}

func NewMarkNotificationAsReadController(useCase *application.MarkNotificationAsReadUseCase) *MarkNotificationAsReadController {
	return &MarkNotificationAsReadController{useCase: useCase}
}

func (c *MarkNotificationAsReadController) Handle(w http.ResponseWriter, r *http.Request) {
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

	updated, err := c.useCase.Execute(r.Context(), id, claims.UserID)
	if err != nil {
		if errors.Is(err, repositories.ErrNotificationNotFound) {
			httpresponse.WriteError(w, http.StatusNotFound, err.Error())
			return
		}
		httpresponse.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	httpresponse.WriteJSON(w, http.StatusOK, updated)
}
