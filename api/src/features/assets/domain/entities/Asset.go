package entities

import "time"

type Asset struct {
	ID              string
	UserID          string
	Name            string
	Category        string
	Brand           string
	PurchaseValue   *float64
	Condition       string
	PurchaseDate    *time.Time
	StoreOrigin     string
	Notes           string
	BlockchainTxID  string
	BlockchainHash  string
	CreatedAt       time.Time
	IsForSale       bool
	SalePrice       *float64
	SaleDescription string
	Size            string
	SellerName      string
	SellerAvatarURL string
}
