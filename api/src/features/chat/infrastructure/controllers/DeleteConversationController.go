package controllers

import (
	"net/http"

	"vault/src/core/httpresponse"
	"vault/src/core/security"
	"vault/src/features/chat/application"
)

type DeleteConversationController struct {
	useCase *application.DeleteConversationUseCase
}

func NewDeleteConversationController(useCase *application.DeleteConversationUseCase) *DeleteConversationController {
	return &DeleteConversationController{useCase: useCase}
}

func (c *DeleteConversationController) Handle(w http.ResponseWriter, r *http.Request) {
	claims, ok := security.ClaimsFromContext(r.Context())
	if !ok {
		httpresponse.WriteError(w, http.StatusUnauthorized, "no autenticado")
		return
	}

	otherID := r.PathValue("id")

	if err := c.useCase.Execute(r.Context(), claims.UserID, otherID); err != nil {
		httpresponse.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
