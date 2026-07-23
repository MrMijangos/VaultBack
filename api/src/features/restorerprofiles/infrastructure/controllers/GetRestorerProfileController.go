package controllers

import (
	"errors"
	"net/http"

	"vault/src/core/httpresponse"
	"vault/src/features/restorerprofiles/application"
	"vault/src/features/restorerprofiles/domain/repositories"
)

type GetRestorerProfileController struct {
	useCase *application.GetRestorerProfileUseCase
}

func NewGetRestorerProfileController(useCase *application.GetRestorerProfileUseCase) *GetRestorerProfileController {
	return &GetRestorerProfileController{useCase: useCase}
}

func (c *GetRestorerProfileController) Handle(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userId")

	profile, err := c.useCase.Execute(r.Context(), userID)
	if err != nil {
		if errors.Is(err, repositories.ErrProfileNotFound) {
			httpresponse.WriteError(w, http.StatusNotFound, err.Error())
			return
		}
		httpresponse.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	httpresponse.WriteJSON(w, http.StatusOK, profile)
}
