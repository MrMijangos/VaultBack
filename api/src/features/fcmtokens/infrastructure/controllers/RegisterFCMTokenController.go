package controllers

import (
	"encoding/json"
	"net/http"

	"vault/src/core/httpresponse"
	"vault/src/core/security"
	"vault/src/features/fcmtokens/application"
	"vault/src/features/fcmtokens/domain/dto/request"
)

type RegisterFCMTokenController struct {
	useCase *application.RegisterFCMTokenUseCase
}

func NewRegisterFCMTokenController(useCase *application.RegisterFCMTokenUseCase) *RegisterFCMTokenController {
	return &RegisterFCMTokenController{useCase: useCase}
}

func (c *RegisterFCMTokenController) Handle(w http.ResponseWriter, r *http.Request) {
	claims, ok := security.ClaimsFromContext(r.Context())
	if !ok {
		httpresponse.WriteError(w, http.StatusUnauthorized, "no autenticado")
		return
	}

	var req request.RegisterFCMTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, "cuerpo de la peticion invalido")
		return
	}

	if err := c.useCase.Execute(r.Context(), claims.UserID, req); err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	httpresponse.WriteJSON(w, http.StatusOK, map[string]string{"message": "token registered"})
}
