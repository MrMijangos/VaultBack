package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"vault/src/core/httpresponse"
	"vault/src/core/security"
	"vault/src/features/recommendations/application"
	"vault/src/features/recommendations/infrastructure/adapters"
)

type GetRecommendationsController struct {
	useCase *application.GetRecommendationsUseCase
}

func NewGetRecommendationsController(useCase *application.GetRecommendationsUseCase) *GetRecommendationsController {
	return &GetRecommendationsController{useCase: useCase}
}

// Handle atiende GET /api/v1/recommendations?top_n=3
// top_n es opcional, default 3, para no obligar al front a mandarlo.
func (c *GetRecommendationsController) Handle(w http.ResponseWriter, r *http.Request) {
	claims, ok := security.ClaimsFromContext(r.Context())
	if !ok {
		httpresponse.WriteError(w, http.StatusUnauthorized, "no autenticado")
		return
	}

	topN := 3
	if raw := r.URL.Query().Get("top_n"); raw != "" {
		if parsed, err := strconv.Atoi(raw); err == nil && parsed > 0 {
			topN = parsed
		}
	}

	result, err := c.useCase.Execute(r.Context(), claims.UserID, topN)
	if err != nil {
		if errors.Is(err, adapters.ErrUnavailable) {
			httpresponse.WriteError(w, http.StatusServiceUnavailable, err.Error())
			return
		}
		httpresponse.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	httpresponse.WriteJSON(w, http.StatusOK, result)
}
