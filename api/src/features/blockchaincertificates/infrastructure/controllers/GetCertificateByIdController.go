package controllers

import (
	"errors"
	"net/http"

	"github.com/google/uuid"

	"vault/src/core/httpresponse"
	"vault/src/features/blockchaincertificates/application"
	"vault/src/features/blockchaincertificates/domain/repositories"
)

type GetCertificateByIdController struct {
	useCase *application.GetCertificateByIdUseCase
}

func NewGetCertificateByIdController(useCase *application.GetCertificateByIdUseCase) *GetCertificateByIdController {
	return &GetCertificateByIdController{useCase: useCase}
}

func (c *GetCertificateByIdController) Handle(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if _, err := uuid.Parse(id); err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, "id invalido")
		return
	}

	cert, err := c.useCase.Execute(r.Context(), id)
	if err != nil {
		if errors.Is(err, repositories.ErrCertificateNotFound) {
			httpresponse.WriteError(w, http.StatusNotFound, err.Error())
			return
		}
		httpresponse.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	httpresponse.WriteJSON(w, http.StatusOK, cert)
}
