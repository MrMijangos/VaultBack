package response

import (
	"time"

	"vault/src/features/maintenancelogs/domain/entities"
)

type MaintenanceLogResponse struct {
	ID             string    `json:"id"`
	AssetID        string    `json:"asset_id"`
	ProviderID     *string   `json:"provider_id"`
	Type           string    `json:"type"`
	Subtype        string    `json:"subtype"`
	Cost           *float64  `json:"cost"`
	PerformedAt    *string   `json:"performed_at"`
	Notes          string    `json:"notes"`
	BlockchainTxID string    `json:"blockchain_tx_id"`
	CreatedAt      time.Time `json:"created_at"`
}

func FromEntity(l entities.MaintenanceLog) MaintenanceLogResponse {
	var performedAt *string
	if l.PerformedAt != nil {
		formatted := l.PerformedAt.Format("2006-01-02")
		performedAt = &formatted
	}

	return MaintenanceLogResponse{
		ID:             l.ID,
		AssetID:        l.AssetID,
		ProviderID:     l.ProviderID,
		Type:           l.Type,
		Subtype:        l.Subtype,
		Cost:           l.Cost,
		PerformedAt:    performedAt,
		Notes:          l.Notes,
		BlockchainTxID: l.BlockchainTxID,
		CreatedAt:      l.CreatedAt,
	}
}

func FromEntities(list []entities.MaintenanceLog) []MaintenanceLogResponse {
	out := make([]MaintenanceLogResponse, 0, len(list))
	for _, l := range list {
		out = append(out, FromEntity(l))
	}
	return out
}
