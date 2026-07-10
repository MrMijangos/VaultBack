package controllers

import (
	"net/http"

	"vault/src/core/httpresponse"
	"vault/src/features/posts/application"
)

type GetAllPostsController struct {
	useCase *application.GetAllPostsUseCase
}

func NewGetAllPostsController(useCase *application.GetAllPostsUseCase) *GetAllPostsController {
	return &GetAllPostsController{useCase: useCase}
}

func (c *GetAllPostsController) Handle(w http.ResponseWriter, r *http.Request) {
	list, err := c.useCase.Execute(r.Context())
	if err != nil {
		httpresponse.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httpresponse.WriteJSON(w, http.StatusOK, list)
}
