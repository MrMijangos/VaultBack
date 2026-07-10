package controllers

import (
	"errors"
	"net/http"

	"github.com/google/uuid"

	"vault/src/core/httpresponse"
	"vault/src/features/posts/application"
	"vault/src/features/posts/domain/repositories"
)

type GetPostByIdController struct {
	useCase *application.GetPostByIdUseCase
}

func NewGetPostByIdController(useCase *application.GetPostByIdUseCase) *GetPostByIdController {
	return &GetPostByIdController{useCase: useCase}
}

func (c *GetPostByIdController) Handle(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if _, err := uuid.Parse(id); err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, "id invalido")
		return
	}

	p, err := c.useCase.Execute(r.Context(), id)
	if err != nil {
		if errors.Is(err, repositories.ErrPostNotFound) {
			httpresponse.WriteError(w, http.StatusNotFound, err.Error())
			return
		}
		httpresponse.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	httpresponse.WriteJSON(w, http.StatusOK, p)
}
