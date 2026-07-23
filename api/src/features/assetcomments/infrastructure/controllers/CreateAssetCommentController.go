package controllers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"

	"vault/src/core/httpresponse"
	"vault/src/core/moderation"
	"vault/src/core/security"
	"vault/src/features/assetcomments/application"
	"vault/src/features/assetcomments/domain/dto/request"
)

type CreateAssetCommentController struct {
	useCase *application.CreateAssetCommentUseCase
}

func NewCreateAssetCommentController(useCase *application.CreateAssetCommentUseCase) *CreateAssetCommentController {
	return &CreateAssetCommentController{useCase: useCase}
}

func (c *CreateAssetCommentController) Handle(w http.ResponseWriter, r *http.Request) {
	claims, ok := security.ClaimsFromContext(r.Context())
	if !ok {
		httpresponse.WriteError(w, http.StatusUnauthorized, "no autenticado")
		return
	}

	assetID := r.PathValue("id")
	if _, err := uuid.Parse(assetID); err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, "id invalido")
		return
	}

	var req request.CreateAssetCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, "cuerpo de la peticion invalido")
		return
	}

	created, err := c.useCase.Execute(r.Context(), assetID, claims.UserID, req)
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
