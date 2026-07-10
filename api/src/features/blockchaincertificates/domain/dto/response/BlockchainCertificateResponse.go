package response

import (
	"time"

	"vault/src/features/blockchaincertificates/domain/entities"
)

type BlockchainCertificateResponse struct {
	ID          string    `json:"id"`
	AssetID     string    `json:"asset_id"`
	OwnerID     string    `json:"owner_id"`
	TxID        string    `json:"tx_id"`
	AssetHash   string    `json:"asset_hash"`
	Action      string    `json:"action"`
	Network     string    `json:"network"`
	ConfirmedAt time.Time `json:"confirmed_at"`
}

func FromEntity(c entities.BlockchainCertificate) BlockchainCertificateResponse {
	return BlockchainCertificateResponse{
		ID:          c.ID,
		AssetID:     c.AssetID,
		OwnerID:     c.OwnerID,
		TxID:        c.TxID,
		AssetHash:   c.AssetHash,
		Action:      c.Action,
		Network:     c.Network,
		ConfirmedAt: c.ConfirmedAt,
	}
}

func FromEntities(list []entities.BlockchainCertificate) []BlockchainCertificateResponse {
	out := make([]BlockchainCertificateResponse, 0, len(list))
	for _, c := range list {
		out = append(out, FromEntity(c))
	}
	return out
}
