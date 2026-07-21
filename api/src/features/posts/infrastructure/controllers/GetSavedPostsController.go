package controllers

import (
	"net/http"

	"vault/src/core/httpresponse"
	"vault/src/core/security"
	"vault/src/features/posts/application"
)

type GetSavedPostsController struct {
	useCase *application.GetSavedPostsUseCase
}

func NewGetSavedPostsController(useCase *application.GetSavedPostsUseCase) *GetSavedPostsController {
	return &GetSavedPostsController{useCase: useCase}
}

func (c *GetSavedPostsController) Handle(w http.ResponseWriter, r *http.Request) {
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
