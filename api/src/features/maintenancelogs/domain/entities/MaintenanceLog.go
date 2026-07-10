package entities

import "time"

type MaintenanceLog struct {
	ID             string
	AssetID        string
	ProviderID     *string
	Type           string
	Subtype        string
	Cost           *float64
	PerformedAt    *time.Time
	Notes          string
	BlockchainTxID string
	CreatedAt      time.Time
}
