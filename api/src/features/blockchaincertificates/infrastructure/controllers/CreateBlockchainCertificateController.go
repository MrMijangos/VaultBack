package controllers

import (
	"encoding/json"
	"errors"
	"net/http"

	"vault/src/core/httpresponse"
	"vault/src/core/security"
	"vault/src/features/blockchaincertificates/application"
	"vault/src/features/blockchaincertificates/domain/dto/request"
	"vault/src/features/blockchaincertificates/domain/repositories"
)

type CreateBlockchainCertificateController struct {
	useCase *application.CreateBlockchainCertificateUseCase
}

func NewCreateBlockchainCertificateController(useCase *application.CreateBlockchainCertificateUseCase) *CreateBlockchainCertificateController {
	return &CreateBlockchainCertificateController{useCase: useCase}
}

func (c *CreateBlockchainCertificateController) Handle(w http.ResponseWriter, r *http.Request) {
	claims, ok := security.ClaimsFromContext(r.Context())
	if !ok {
		httpresponse.WriteError(w, http.StatusUnauthorized, "no autenticado")
		return
	}

	var req request.CreateBlockchainCertificateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, "cuerpo de la peticion invalido")
		return
	}

	created, err := c.useCase.Execute(r.Context(), claims.UserID, req)
	if err != nil {
		if errors.Is(err, repositories.ErrTxIDAlreadyExists) {
			httpresponse.WriteError(w, http.StatusConflict, err.Error())
			return
		}
		httpresponse.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	httpresponse.WriteJSON(w, http.StatusCreated, created)
}
