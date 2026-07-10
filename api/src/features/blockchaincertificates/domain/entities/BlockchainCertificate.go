package entities

import "time"

type BlockchainCertificate struct {
	ID          string
	AssetID     string
	OwnerID     string
	TxID        string
	AssetHash   string
	Action      string
	Network     string
	ConfirmedAt time.Time
}
