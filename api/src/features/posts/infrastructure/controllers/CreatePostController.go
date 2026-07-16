package controllers

import (
	"encoding/json"
	"errors"
	"net/http"

	"vault/src/core/httpresponse"
	"vault/src/core/moderation"
	"vault/src/core/security"
	"vault/src/features/posts/application"
	"vault/src/features/posts/domain/dto/request"
)

type CreatePostController struct {
	useCase *application.CreatePostUseCase
}

func NewCreatePostController(useCase *application.CreatePostUseCase) *CreatePostController {
	return &CreatePostController{useCase: useCase}
}

func (c *CreatePostController) Handle(w http.ResponseWriter, r *http.Request) {
	claims, ok := security.ClaimsFromContext(r.Context())
	if !ok {
		httpresponse.WriteError(w, http.StatusUnauthorized, "no autenticado")
		return
	}

	var req request.CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, "cuerpo de la peticion invalido")
		return
	}

	created, err := c.useCase.Execute(r.Context(), claims.UserID, req)
	if err != nil {
		switch {
		case errors.Is(err, moderation.ErrToxicContent):
			httpresponse.WriteError(w, http.StatusUnprocessableEntity, err.Error())
		case errors.Is(err, moderation.ErrUnavailable):
			httpresponse.WriteError(w, http.StatusServiceUnavailable, err.Error())
		default:
			httpresponse.WriteError(w, http.StatusBadRequest, err.Error())
		}
		return
	}

	httpresponse.WriteJSON(w, http.StatusCreated, created)
}
