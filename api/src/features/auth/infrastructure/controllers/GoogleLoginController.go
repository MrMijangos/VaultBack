package controllers

import (
	"encoding/json"
	"net/http"

	"vault/src/core/httpresponse"
	"vault/src/core/security"
	"vault/src/features/auth/application"
	"vault/src/features/auth/domain/dto/request"
)

type GoogleLoginController struct {
	useCase      *application.GoogleLoginUseCase
	cookieSecure bool
}

func NewGoogleLoginController(useCase *application.GoogleLoginUseCase, cookieSecure bool) *GoogleLoginController {
	return &GoogleLoginController{useCase: useCase, cookieSecure: cookieSecure}
}

func (c *GoogleLoginController) Handle(w http.ResponseWriter, r *http.Request) {
	var req request.GoogleLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, "cuerpo de la peticion invalido")
		return
	}

	user, token, err := c.useCase.Execute(r.Context(), req)
	if err != nil {
		httpresponse.WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}

	security.SetAuthCookie(w, token, c.cookieSecure)
	user.Token = token
	httpresponse.WriteJSON(w, http.StatusOK, user)
}
