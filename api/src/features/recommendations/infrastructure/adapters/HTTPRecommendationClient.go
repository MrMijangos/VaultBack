package adapters

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"vault/src/features/recommendations/domain/dto/response"
)

// ErrUnavailable se devuelve cuando no se pudo consultar vault-ai-service
// (caído, timeout, respuesta inválida) o cuando el servicio contesta que
// el usuario aún no tiene suficiente historial para recomendar (422).
var ErrUnavailable = errors.New("no se pudieron obtener recomendaciones en este momento")

// HTTPRecommendationClient llama a POST /api/v1/ml/recommend-items de
// vault-ai-service (mismo servicio que moderation.Client usa para NLP,
// por eso reutiliza la misma baseURL: cfg.NLPServiceURL).
type HTTPRecommendationClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewHTTPRecommendationClient(baseURL string) *HTTPRecommendationClient {
	return &HTTPRecommendationClient{
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

type recommendItemsRequest struct {
	UserID string `json:"user_id"`
}

func (c *HTTPRecommendationClient) GetRecommendedItems(ctx context.Context, userID string, topN int) (response.RecommendationsResponse, error) {
	body, err := json.Marshal(recommendItemsRequest{UserID: userID})
	if err != nil {
		return response.RecommendationsResponse{}, fmt.Errorf("%w: no se pudo serializar la solicitud: %v", ErrUnavailable, err)
	}

	url := c.baseURL + "/api/v1/ml/recommend-items?top_n=" + strconv.Itoa(topN)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return response.RecommendationsResponse{}, fmt.Errorf("%w: no se pudo crear la solicitud: %v", ErrUnavailable, err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return response.RecommendationsResponse{}, fmt.Errorf("%w: %v", ErrUnavailable, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return response.RecommendationsResponse{}, fmt.Errorf("%w: el servicio de recomendaciones devolvio %d", ErrUnavailable, resp.StatusCode)
	}

	var out response.RecommendationsResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return response.RecommendationsResponse{}, fmt.Errorf("%w: respuesta invalida: %v", ErrUnavailable, err)
	}

	return out, nil
}
