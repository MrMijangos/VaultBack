package controllers

import (
	"net/http"

	"github.com/google/uuid"

	"vault/src/core/httpresponse"
	"vault/src/core/security"
	"vault/src/features/reviews/application"
)

type LikeReviewController struct {
	useCase *application.LikeReviewUseCase
}

func NewLikeReviewController(useCase *application.LikeReviewUseCase) *LikeReviewController {
	return &LikeReviewController{useCase: useCase}
}

func (c *LikeReviewController) Handle(w http.ResponseWriter, r *http.Request) {
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

	if err := c.useCase.Execute(r.Context(), id, claims.UserID); err != nil {
		httpresponse.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
