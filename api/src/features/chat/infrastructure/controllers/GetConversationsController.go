package controllers

import (
	"net/http"

	"vault/src/core/httpresponse"
	"vault/src/core/security"
	"vault/src/features/chat/application"
)

type GetConversationsController struct {
	useCase *application.GetConversationsUseCase
}

func NewGetConversationsController(useCase *application.GetConversationsUseCase) *GetConversationsController {
	return &GetConversationsController{useCase: useCase}
}

func (c *GetConversationsController) Handle(w http.ResponseWriter, r *http.Request) {
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
