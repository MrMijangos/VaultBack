package controllers

import (
	"encoding/json"
	"net/http"

	"vault/src/core/httpresponse"
	"vault/src/core/security"
	"vault/src/features/fcmtokens/application"
	"vault/src/features/fcmtokens/domain/dto/request"
)

type DeleteFCMTokenController struct {
	useCase *application.DeleteFCMTokenUseCase
}

func NewDeleteFCMTokenController(useCase *application.DeleteFCMTokenUseCase) *DeleteFCMTokenController {
	return &DeleteFCMTokenController{useCase: useCase}
}

func (c *DeleteFCMTokenController) Handle(w http.ResponseWriter, r *http.Request) {
	_, ok := security.ClaimsFromContext(r.Context())
	if !ok {
		httpresponse.WriteError(w, http.StatusUnauthorized, "no autenticado")
		return
	}

	var req request.DeleteFCMTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, "cuerpo de la peticion invalido")
		return
	}

	if err := c.useCase.Execute(r.Context(), req); err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}
