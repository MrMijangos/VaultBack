package controllers

import (
	"net/http"

	"vault/src/core/httpresponse"
	"vault/src/features/reviews/application"
)

type GetProviderRatingController struct {
	useCase *application.GetProviderRatingUseCase
}

func NewGetProviderRatingController(useCase *application.GetProviderRatingUseCase) *GetProviderRatingController {
	return &GetProviderRatingController{useCase: useCase}
}

func (c *GetProviderRatingController) Handle(w http.ResponseWriter, r *http.Request) {
	providerID := r.URL.Query().Get("provider_id")
	if providerID == "" {
		httpresponse.WriteError(w, http.StatusBadRequest, "el parametro provider_id es obligatorio")
		return
	}

	rating, err := c.useCase.Execute(r.Context(), providerID)
	if err != nil {
		httpresponse.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httpresponse.WriteJSON(w, http.StatusOK, rating)
}
