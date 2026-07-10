package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	"vault/src/core/httpresponse"
	"vault/src/core/security"
	"vault/src/features/comments/application"
	"vault/src/features/comments/domain/dto/request"
)

type CreateCommentController struct {
	useCase *application.CreateCommentUseCase
}

func NewCreateCommentController(useCase *application.CreateCommentUseCase) *CreateCommentController {
	return &CreateCommentController{useCase: useCase}
}

func (c *CreateCommentController) Handle(w http.ResponseWriter, r *http.Request) {
	claims, ok := security.ClaimsFromContext(r.Context())
	if !ok {
		httpresponse.WriteError(w, http.StatusUnauthorized, "no autenticado")
		return
	}

	postID := r.PathValue("id")
	if _, err := uuid.Parse(postID); err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, "id invalido")
		return
	}

	var req request.CreateCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, "cuerpo de la peticion invalido")
		return
	}

	created, err := c.useCase.Execute(r.Context(), postID, claims.UserID, req)
	if err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	httpresponse.WriteJSON(w, http.StatusCreated, created)
}
