package controllers

import (
	"net/http"

	"vault/src/core/httpresponse"
	"vault/src/features/businesses/application"
)

type GetAllBusinessesController struct {
	useCase *application.GetAllBusinessesUseCase
}

func NewGetAllBusinessesController(useCase *application.GetAllBusinessesUseCase) *GetAllBusinessesController {
	return &GetAllBusinessesController{useCase: useCase}
}

func (c *GetAllBusinessesController) Handle(w http.ResponseWriter, r *http.Request) {
	list, err := c.useCase.Execute(r.Context())
	if err != nil {
		httpresponse.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httpresponse.WriteJSON(w, http.StatusOK, list)
}
