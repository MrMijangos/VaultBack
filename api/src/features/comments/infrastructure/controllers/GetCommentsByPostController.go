package controllers

import (
	"net/http"

	"github.com/google/uuid"

	"vault/src/core/httpresponse"
	"vault/src/features/comments/application"
)

type GetCommentsByPostController struct {
	useCase *application.GetCommentsByPostUseCase
}

func NewGetCommentsByPostController(useCase *application.GetCommentsByPostUseCase) *GetCommentsByPostController {
	return &GetCommentsByPostController{useCase: useCase}
}

func (c *GetCommentsByPostController) Handle(w http.ResponseWriter, r *http.Request) {
	postID := r.PathValue("id")
	if _, err := uuid.Parse(postID); err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, "id invalido")
		return
	}

	list, err := c.useCase.Execute(r.Context(), postID)
	if err != nil {
		httpresponse.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httpresponse.WriteJSON(w, http.StatusOK, list)
}
