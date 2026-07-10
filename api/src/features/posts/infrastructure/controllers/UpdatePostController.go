package controllers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"

	"vault/src/core/httpresponse"
	"vault/src/core/security"
	"vault/src/features/posts/application"
	"vault/src/features/posts/domain/dto/request"
	"vault/src/features/posts/domain/repositories"
)

type UpdatePostController struct {
	useCase *application.UpdatePostUseCase
}

func NewUpdatePostController(useCase *application.UpdatePostUseCase) *UpdatePostController {
	return &UpdatePostController{useCase: useCase}
}

func (c *UpdatePostController) Handle(w http.ResponseWriter, r *http.Request) {
	claims, ok := security.ClaimsFromContext(r.Context())
	if !ok {
		httpresponse.WriteError(w, http.StatusUnauthorized, "no autenticado")
		return
	}

	id := r.PathValue("id")
	if _, err := uuid.Parse(id); err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, "id invalido")
		return
	}

	var req request.UpdatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, "cuerpo de la peticion invalido")
		return
	}

	updated, err := c.useCase.Execute(r.Context(), id, claims.UserID, req)
	if err != nil {
		if errors.Is(err, repositories.ErrPostNotFound) {
			httpresponse.WriteError(w, http.StatusNotFound, err.Error())
			return
		}
		httpresponse.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	httpresponse.WriteJSON(w, http.StatusOK, updated)
}
