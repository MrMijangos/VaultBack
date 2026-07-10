package controllers

import (
	"net/http"

	"vault/src/core/httpresponse"
	"vault/src/core/security"
	"vault/src/features/notifications/application"
)

type GetMyNotificationsController struct {
	useCase *application.GetMyNotificationsUseCase
}

func NewGetMyNotificationsController(useCase *application.GetMyNotificationsUseCase) *GetMyNotificationsController {
	return &GetMyNotificationsController{useCase: useCase}
}

func (c *GetMyNotificationsController) Handle(w http.ResponseWriter, r *http.Request) {
	claims, ok := security.ClaimsFromContext(r.Context())
	if !ok {
		httpresponse.WriteError(w, http.StatusUnauthorized, "no autenticado")
		return
	}

	list, err := c.useCase.Execute(r.Context(), claims.UserID)
	if err != nil {
		httpresponse.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httpresponse.WriteJSON(w, http.StatusOK, list)
}
