package controllers

import (
	"net/http"

	"vault/src/core/httpresponse"
	"vault/src/features/restorerprofiles/application"
)

type ListRestorerProfilesController struct {
	useCase *application.ListRestorerProfilesUseCase
}

func NewListRestorerProfilesController(useCase *application.ListRestorerProfilesUseCase) *ListRestorerProfilesController {
	return &ListRestorerProfilesController{useCase: useCase}
}

func (c *ListRestorerProfilesController) Handle(w http.ResponseWriter, r *http.Request) {
	list, err := c.useCase.Execute(r.Context())
	if err != nil {
		httpresponse.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	httpresponse.WriteJSON(w, http.StatusOK, list)
}
