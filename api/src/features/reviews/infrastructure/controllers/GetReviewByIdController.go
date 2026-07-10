package controllers

import (
	"errors"
	"net/http"

	"github.com/google/uuid"

	"vault/src/core/httpresponse"
	"vault/src/features/reviews/application"
	"vault/src/features/reviews/domain/repositories"
)

type GetReviewByIdController struct {
	useCase *application.GetReviewByIdUseCase
}

func NewGetReviewByIdController(useCase *application.GetReviewByIdUseCase) *GetReviewByIdController {
	return &GetReviewByIdController{useCase: useCase}
}

func (c *GetReviewByIdController) Handle(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if _, err := uuid.Parse(id); err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, "id invalido")
		return
	}

	rv, err := c.useCase.Execute(r.Context(), id)
	if err != nil {
		if errors.Is(err, repositories.ErrReviewNotFound) {
			httpresponse.WriteError(w, http.StatusNotFound, err.Error())
			return
		}
		httpresponse.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	httpresponse.WriteJSON(w, http.StatusOK, rv)
}
