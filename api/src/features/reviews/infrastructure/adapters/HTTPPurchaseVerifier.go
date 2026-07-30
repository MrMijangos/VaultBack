package adapters

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// HTTPPurchaseVerifier llama a GET /api/v1/orders/has-purchased en
// payment/ (ruta pública, no requiere token -- ver
// payment/src/features/orders/infrastructure/router/router.go), mismo
// patrón que payment/ usa para consultar el precio real de un activo en
// api/ (HTTPAssetPriceProvider).
type HTTPPurchaseVerifier struct {
	baseURL string
	client  *http.Client
}

func NewHTTPPurchaseVerifier(baseURL string) *HTTPPurchaseVerifier {
	return &HTTPPurchaseVerifier{baseURL: baseURL, client: &http.Client{Timeout: 10 * time.Second}}
}

type hasPurchasedResponse struct {
	Purchased bool `json:"purchased"`
}

func (p *HTTPPurchaseVerifier) HasPurchased(ctx context.Context, buyerID string, providerID string) (bool, error) {
	query := url.Values{"buyer_id": {buyerID}, "seller_id": {providerID}}
	endpoint := fmt.Sprintf("%s/api/v1/orders/has-purchased?%s", p.baseURL, query.Encode())

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return false, fmt.Errorf("no se pudo armar la consulta de compra: %w", err)
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return false, fmt.Errorf("no se pudo verificar la compra: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("payment respondió %d al verificar la compra", resp.StatusCode)
	}

	var out hasPurchasedResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return false, fmt.Errorf("respuesta invalida al verificar la compra: %w", err)
	}
	return out.Purchased, nil
}
