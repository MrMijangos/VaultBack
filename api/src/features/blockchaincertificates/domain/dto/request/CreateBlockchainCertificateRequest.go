package request

import "errors"

var allowedActions = map[string]bool{
	"REGISTERED":  true,
	"MAINTAINED":  true,
	"RESTORED":    true,
	"TRANSFERRED": true,
}

var allowedNetworks = map[string]bool{
	"testnet": true,
	"mainnet": true,
}

type CreateBlockchainCertificateRequest struct {
	AssetID   string `json:"asset_id"`
	TxID      string `json:"tx_id"`
	AssetHash string `json:"asset_hash"`
	Action    string `json:"action"`
	Network   string `json:"network"`
}

func (r *CreateBlockchainCertificateRequest) Validate() error {
	if r.AssetID == "" {
		return errors.New("el producto es obligatorio")
	}
	if r.TxID == "" {
		return errors.New("el tx_id es obligatorio")
	}
	if r.AssetHash == "" {
		return errors.New("el asset_hash es obligatorio")
	}
	if !allowedActions[r.Action] {
		return errors.New("la accion no es valida")
	}
	if r.Network == "" {
		r.Network = "testnet"
	}
	if !allowedNetworks[r.Network] {
		return errors.New("la red no es valida")
	}
	return nil
}
