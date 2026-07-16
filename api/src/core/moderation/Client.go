package moderation

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// ErrToxicContent se devuelve cuando vault-ai-service marca el contenido
// como tóxico -- la publicación debe rechazarse, no guardarse.
var ErrToxicContent = errors.New("tu contenido es demasiado ofensivo para publicarse")

// ErrUnavailable se devuelve cuando no se pudo consultar el servicio de
// moderación (caído, timeout, respuesta inválida). Con la política de
// "bloquear si no responde", el llamador debe tratar esto como un rechazo
// de la publicación, no como "dejar pasar".
var ErrUnavailable = errors.New("no se pudo verificar el contenido, intenta de nuevo en un momento")

type Result struct {
	SentimentScore float64
	SentimentLabel string
	ToxicityScore  float64
	IsToxic        bool
}

// Client llama sincrónicamente a POST /api/v1/nlp/analyze de
// vault-ai-service antes de guardar un post/comentario/reseña, para
// decidir si se publica o se rechaza -- en vez del flujo anterior de
// publicar primero y corregir is_visible después vía RabbitMQ.
type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

type analyzeRequest struct {
	Text       string `json:"text"`
	SourceID   string `json:"source_id"`
	SourceType string `json:"source_type"`
}

type analyzeResponse struct {
	SentimentScore float64 `json:"sentiment_score"`
	SentimentLabel string  `json:"sentiment_label"`
	ToxicityScore  float64 `json:"toxicity_score"`
	IsToxic        bool    `json:"is_toxic"`
}

// Analyze corre el análisis de NLP sobre text. sourceID/sourceType solo
// correlacionan la solicitud en logs de vault-ai-service -- el análisis en
// sí no depende de ellos. No devuelve ErrToxicContent: eso lo decide el
// llamador según Result.IsToxic, para poder loguear el score antes de
// rechazar.
func (c *Client) Analyze(ctx context.Context, sourceID string, sourceType string, text string) (Result, error) {
	body, err := json.Marshal(analyzeRequest{Text: text, SourceID: sourceID, SourceType: sourceType})
	if err != nil {
		return Result{}, fmt.Errorf("%w: no se pudo serializar la solicitud: %v", ErrUnavailable, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/api/v1/nlp/analyze", bytes.NewReader(body))
	if err != nil {
		return Result{}, fmt.Errorf("%w: no se pudo crear la solicitud: %v", ErrUnavailable, err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return Result{}, fmt.Errorf("%w: %v", ErrUnavailable, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Result{}, fmt.Errorf("%w: el servicio de analisis devolvio %d", ErrUnavailable, resp.StatusCode)
	}

	var out analyzeResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return Result{}, fmt.Errorf("%w: respuesta invalida: %v", ErrUnavailable, err)
	}

	return Result{
		SentimentScore: out.SentimentScore,
		SentimentLabel: out.SentimentLabel,
		ToxicityScore:  out.ToxicityScore,
		IsToxic:        out.IsToxic,
	}, nil
}
