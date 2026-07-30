package adapters

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"time"
)

var ErrAssetNotFound = errors.New("el activo no existe")

type HTTPAssetPriceProvider struct {
	baseURL string
	client  *http.Client
}

func NewHTTPAssetPriceProvider(baseURL string) *HTTPAssetPriceProvider {
	return &HTTPAssetPriceProvider{baseURL: baseURL, client: &http.Client{Timeout: 10 * time.Second}}
}

// assetResponse solo lee los campos que hacen falta -- debe coincidir con
// (un subconjunto de) AssetResponse en api/src/features/assets/domain/dto/response.
type assetResponse struct {
	IsForSale bool     `json:"is_for_sale"`
	SalePrice *float64 `json:"sale_price"`
}

// GetSalePriceCents llama a GET /api/v1/assets/{id} en api/ (ruta pública,
// no requiere token) y convierte sale_price (pesos, ej. 199.99) a centavos.
func (p *HTTPAssetPriceProvider) GetSalePriceCents(ctx context.Context, assetID string) (int64, bool, error) {
	url := fmt.Sprintf("%s/api/v1/assets/%s", p.baseURL, assetID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, false, fmt.Errorf("no se pudo armar la consulta del activo: %w", err)
	}

	resp, err := p.client.Do(req)
	if err != nil {
		return 0, false, fmt.Errorf("no se pudo consultar el precio del activo: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return 0, false, ErrAssetNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return 0, false, fmt.Errorf("api respondió %d al consultar el activo %s", resp.StatusCode, assetID)
	}

	var asset assetResponse
	if err := json.NewDecoder(resp.Body).Decode(&asset); err != nil {
		return 0, false, fmt.Errorf("respuesta invalida al consultar el activo: %w", err)
	}
	if asset.SalePrice == nil {
		return 0, asset.IsForSale, nil
	}
	return int64(math.Round(*asset.SalePrice * 100)), asset.IsForSale, nil
}
