package controllers

import (
	"encoding/json"
	"errors"
	"net/http"

	"vault/src/core/httpresponse"
	"vault/src/core/security"
	"vault/src/features/auth/application"
	"vault/src/features/auth/domain/dto/request"
)

type LoginController struct {
	useCase      *application.LoginUseCase
	cookieSecure bool
}

func NewLoginController(useCase *application.LoginUseCase, cookieSecure bool) *LoginController {
	return &LoginController{useCase: useCase, cookieSecure: cookieSecure}
}

func (c *LoginController) Handle(w http.ResponseWriter, r *http.Request) {
	var req request.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, "cuerpo de la peticion invalido")
		return
	}

	user, token, err := c.useCase.Execute(r.Context(), req)
	if err != nil {
		if errors.Is(err, application.ErrInvalidCredentials) {
			httpresponse.WriteError(w, http.StatusUnauthorized, err.Error())
			return
		}
		httpresponse.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	security.SetAuthCookie(w, token, c.cookieSecure)
	httpresponse.WriteJSON(w, http.StatusOK, user)
}
