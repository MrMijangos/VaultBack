package controllers

import (
	"errors"
	"net/http"

	"github.com/google/uuid"

	"vault/src/core/httpresponse"
	"vault/src/core/security"
	"vault/src/features/assetcomments/application"
	"vault/src/features/assetcomments/domain/repositories"
)

type DeleteAssetCommentController struct {
	useCase *application.DeleteAssetCommentUseCase
}

func NewDeleteAssetCommentController(useCase *application.DeleteAssetCommentUseCase) *DeleteAssetCommentController {
	return &DeleteAssetCommentController{useCase: useCase}
}

func (c *DeleteAssetCommentController) Handle(w http.ResponseWriter, r *http.Request) {
	claims, ok := security.ClaimsFromContext(r.Context())
	if !ok {
		httpresponse.WriteError(w, http.StatusUnauthorized, "no autenticado")
		return
	}

	commentID := r.PathValue("commentId")
	if _, err := uuid.Parse(commentID); err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, "id invalido")
		return
	}

	if err := c.useCase.Execute(r.Context(), commentID, claims.UserID); err != nil {
		if errors.Is(err, repositories.ErrAssetCommentNotFound) {
			httpresponse.WriteError(w, http.StatusNotFound, err.Error())
			return
		}
		httpresponse.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
