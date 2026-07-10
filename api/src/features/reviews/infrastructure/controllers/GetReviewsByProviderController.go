package controllers

import (
	"net/http"

	"vault/src/core/httpresponse"
	"vault/src/features/reviews/application"
)

type GetReviewsByProviderController struct {
	useCase *application.GetReviewsByProviderUseCase
}

func NewGetReviewsByProviderController(useCase *application.GetReviewsByProviderUseCase) *GetReviewsByProviderController {
	return &GetReviewsByProviderController{useCase: useCase}
}

func (c *GetReviewsByProviderController) Handle(w http.ResponseWriter, r *http.Request) {
	providerID := r.URL.Query().Get("provider_id")
	if providerID == "" {
		httpresponse.WriteError(w, http.StatusBadRequest, "el parametro provider_id es obligatorio")
		return
	}

	list, err := c.useCase.Execute(r.Context(), providerID)
	if err != nil {
		httpresponse.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httpresponse.WriteJSON(w, http.StatusOK, list)
}
