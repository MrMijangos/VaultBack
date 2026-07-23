package controllers

import (
	"net/http"

	"vault/src/core/httpresponse"
	"vault/src/core/security"
	"vault/src/features/chat/application"
)

type GetConversationMessagesController struct {
	useCase *application.GetConversationMessagesUseCase
}

func NewGetConversationMessagesController(useCase *application.GetConversationMessagesUseCase) *GetConversationMessagesController {
	return &GetConversationMessagesController{useCase: useCase}
}

func (c *GetConversationMessagesController) Handle(w http.ResponseWriter, r *http.Request) {
	claims, ok := security.ClaimsFromContext(r.Context())
	if !ok {
		httpresponse.WriteError(w, http.StatusUnauthorized, "no autenticado")
		return
	}

	otherID := r.PathValue("id")

	list, err := c.useCase.Execute(r.Context(), claims.UserID, otherID)
	if err != nil {
		httpresponse.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	httpresponse.WriteJSON(w, http.StatusOK, list)
}
